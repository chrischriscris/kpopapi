package main

import (
	"context"
	"fmt"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/chrischriscris/kpopapi/internal/db/utils"
	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)

const baseURL = "https://kpopping.com"
const baseDir = "images"

// ===========  Utils ===========

func getImageDimensions(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

func extractLabel(text string) string {
	return strings.Split(text, "\n")[3][3:]
}

func buildPhotoNameFromURL(url string) string {
	p := strings.Split(url, "/")
	n := p[len(p)-1]

	return strings.Split(n, "?")[0]
}

func downloadImage(url string, directory string, photo string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	imgFullPath := fmt.Sprintf("%s/%s.jpg", directory, buildPhotoNameFromURL(photo))
	os.MkdirAll(directory, os.ModePerm)
	out, err := os.Create(imgFullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return imgFullPath, nil
}

// =========== Database logic ===========

func registerAndSaveImage(
	ctx context.Context,
	qtx *repository.Queries,
	url string,
    directory string,
	photo string,
) error {

	// Early return if the image already exists in the database
	_, err := qtx.GetImageByUrl(ctx, url)
	if err == nil {
		return nil
	}

	entityID, isGroup, err := getGroupOrIdolIDFromDB(ctx, *qtx, directory)
	if err != nil {
		return err
	}

	outDir := fmt.Sprintf("%s/%s", baseDir, directory)
	imageID, err := downloadImageAndSaveToDB(ctx, *qtx, url, outDir, photo)
	if err != nil {
        return err
	}

	if isGroup {
		_, err = qtx.AddGroupImage(ctx, repository.AddGroupImageParams{
			ImageID: imageID,
			GroupID: entityID,
		})
	} else {
		_, err = qtx.AddIdolImage(ctx, repository.AddIdolImageParams{
			ImageID: imageID,
			IdolID:  entityID,
		})
	}

	return err
}

// Fetches the ID of a group or idol from the database
// Returns the ID and a true if the entity was a group, false if it was an idol
func getGroupOrIdolIDFromDB(
	ctx context.Context,
	qtx repository.Queries,
	name string,
) (int32, bool, error) {
	group, err := qtx.GetGroupByName(ctx, name)
	if err == nil {
		return group.ID, true, nil
	}

	idol, err := qtx.GetIdolByName(ctx, pgtype.Text{String: name})
	if err == nil {
		return idol.ID, false, nil
	}

	return 0, false, fmt.Errorf("Unable to find group or idol with name %s", name)
}

// Downloads an image from a URL and saves it to the database
// Returns the ID of the image and an error
func downloadImageAndSaveToDB(
	ctx context.Context,
	qtx repository.Queries,
	url string,
	outDir string,
	photo string,
) (int32, error) {
	imgPath, err := downloadImage(url, outDir, photo)
	if err != nil {
		return 0, fmt.Errorf("Unable to download image: %v\n", err)
	}

	width, height, err := getImageDimensions(imgPath)
	if err != nil {
		return 0, fmt.Errorf("Unable to get image dimensions: %v\n", err)
	}

	image, err := qtx.AddImage(ctx, url)
	if err != nil {
		return 0, fmt.Errorf("Unable to add image: %v\n", err)
	}

	_, err = qtx.AddImageMetadata(ctx, repository.AddImageMetadataParams{
		ImageID: image.ID,
		Width:   int32(width),
		Height:  int32(height),
	})

	return image.ID, err
}

// =========== Scraping logic ===========

func getPageLinks() map[string][]string {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting scraping: ", r.URL)
	})

	// links[label] = [link1, link2, ...]
	links := make(map[string][]string)
	c.OnHTML("div.cell", func(e *colly.HTMLElement) {
		artist := extractLabel(e.DOM.Find("figcaption").First().Text())
		link := e.ChildAttr("a[aria-label='picture']", "href")

		links[artist] = append(links[artist], link)
	})

	c.Visit(baseURL + "/kpics")

	return links
}

func downloadImages(links map[string][]string) int {
	acc := 0
	for artist, links := range links {
		acc += downloadArtistImages(artist, links)
	}

	return acc
}

func downloadArtistImages(artist string, links []string) int {
	acc := 0
	for _, link := range links {
		acc += downloadImagesFromLink(artist, link)
	}

	return acc
}

func downloadImagesFromLink(directory string, link string) int {
	c := colly.NewCollector()

	photos := make([]string, 0)
	c.OnHTML("div.justified-gallery a", func(e *colly.HTMLElement) {
		photo := e.Attr("href")

		photos = append(photos, photo)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Found", len(photos), "photos for", directory)
		ctx, conn, err := utils.ConnectDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		defer conn.Close(context.Background())

		tx, qtx, err := utils.BeginTransaction(ctx, conn)
		if err != nil {
			log.Fatalf("Unable to start transaction: %v\n", err)
		}
		defer tx.Rollback(ctx)

		for _, photo := range photos {
			downloadURL := baseURL + photo

            err := registerAndSaveImage(ctx, qtx, downloadURL, directory, photo)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatalf("Unable to commit transaction: %v\n", err)
		}
	})

	c.Visit(baseURL + link)

	return len(photos)
}

func main() {
	links := getPageLinks()
	total := downloadImages(links)

	fmt.Println("Downloaded", total, "images")
}
