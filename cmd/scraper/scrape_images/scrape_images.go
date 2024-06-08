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

func saveImageToDB() {
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

    _, err := qtx.AddImage(ctx, "https://kpopping.com/kpics/2021/12/20211207-0001-0001.jpg")
    if err != nil {
        log.Fatalf("Unable to add image: %v\n", err)
    }

    // imageMetadata, err := qtx.AddImageMetadata(ctx, repository.AddImageMetadataParams{
    //     ImageID: image.ID,
    //     Width:
    //
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
		for _, photo := range photos {
			downloadURL := baseURL + photo

			outDir := fmt.Sprintf("%s/%s", baseDir, directory)
			imgPath, err := downloadImage(downloadURL, outDir, photo)
			if err != nil {
				log.Fatal(err)
			}

			width, height, err := getImageDimensions(imgPath)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Downloaded", imgPath, "with dimensions", width, "x", height)
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
