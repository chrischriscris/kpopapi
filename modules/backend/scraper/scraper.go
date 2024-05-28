package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

const baseURL = "https://kpopping.com"
const baseDir = "images"

func extractLabel(text string) string {
	return strings.Split(text, "\n")[3][3:]
}

func buildPhotoNameFromURL(url string) string {
	p := strings.Split(url, "/")
	n := p[len(p)-1]

	return strings.Split(n, "?")[0]
}

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

			resp, err := http.Get(downloadURL)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

            outDir := fmt.Sprintf("%s/%s", baseDir, directory)
			os.MkdirAll(outDir, os.ModePerm)
			out, err := os.Create(fmt.Sprintf("%s/%s.jpg", directory, buildPhotoNameFromURL(photo)))
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
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
