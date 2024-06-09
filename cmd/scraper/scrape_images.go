package main

import (
	"context"
	"fmt"

	"log"

	"github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/chrischriscris/kpopapi/internal/scraper"
	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)


// =========== Database logic ===========

func registerAndSaveImage(
	ctx context.Context,
	qtx *repository.Queries,
	artist string,
	url string,
) error {

	// Early return if the image already exists in the database
	_, err := qtx.GetImageByUrl(ctx, url)
	if err == nil {
		return fmt.Errorf("Image already exists in the database")
	}

	entityID, isGroup, err := getGroupOrIdolIDFromDB(ctx, *qtx, artist)
	if err != nil {
		return err
	}

	outDir := fmt.Sprintf("%s/%s", scraperutils.BaseDir, artist)
	imageID, err := downloadImageAndSaveToDB(ctx, *qtx, url, outDir)
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

	idol, err := qtx.GetIdolByName(ctx, pgtype.Text{String: name, Valid: true})
	if err == nil {
		return idol.ID, false, nil
	}

	return 0, false, fmt.Errorf("Unable to find group or idol with name '%s'", name)
}

// Downloads an image from a URL and saves it to the database
// Returns the ID of the image and an error
func downloadImageAndSaveToDB(
	ctx context.Context,
	qtx repository.Queries,
	url string,
	outDir string,
) (int32, error) {
	imgPath, err := scraperutils.DownloadImage(url, outDir)
	if err != nil {
		return 0, fmt.Errorf("Unable to download image: %v\n", err)
	}

	width, height, err := scraperutils.GetImageDimensions(imgPath)
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
		artist := scraperutils.ExtractLabel(e.DOM.Find("figcaption").First().Text())
		link := e.ChildAttr("a[aria-label='picture']", "href")

		links[artist] = append(links[artist], link)
	})

	c.Visit(scraperutils.BaseURL + "/kpics")

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

func downloadImagesFromLink(artist string, link string) int {
	c := colly.NewCollector()

	photoURLs := make([]string, 0)
	c.OnHTML("div.justified-gallery a", func(e *colly.HTMLElement) {
		photoPath := e.Attr("href")
		photoURL := scraperutils.BaseURL + photoPath

		photoURLs = append(photoURLs, photoURL)
	})

	n_downloaded := 0
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Found", len(photoURLs), "photos for", artist)
		ctx, conn, err := dbutils.ConnectDB()
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		defer conn.Close(context.Background())

		tx, qtx, err := dbutils.BeginTransaction(ctx, conn)
		if err != nil {
			log.Fatalf("Unable to start transaction: %v\n", err)
		}
		defer tx.Rollback(ctx)

		for _, photoURL := range photoURLs {
			err := registerAndSaveImage(ctx, qtx, artist, photoURL)
			if err != nil {
				if err.Error() != "Image already exists in the database" {
					fmt.Printf("Error: %v\n", err)
				}
				continue
			}
			n_downloaded++
		}

		err = tx.Commit(ctx)
		if err != nil {
			log.Fatalf("Unable to commit transaction: %v\n", err)
		}
	})

	c.Visit(scraperutils.BaseURL + link)

	return n_downloaded
}

func main() {
	links := getPageLinks()
	total := downloadImages(links)

	fmt.Println("Downloaded", total, "images")
}
