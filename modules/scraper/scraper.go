package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly"
)

const baseURL = "https://kpopping.com"

func main() {
	// Scraping logic here
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	links := make([]string, 0)
    first := false

	c.OnHTML("a[aria-label='picture']", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		links = append(links, link)
        if !first {
            first = true
            c.Visit(baseURL + link)
        }
	})

	photos := make([]string, 0)
	// Grab all <a> inside <div> with class 'box pics'
	c.OnHTML("div.justified-gallery a", func(e *colly.HTMLElement) {
		photo := e.Attr("href")

		photos = append(photos, photo)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Links found: ")
		for _, link := range links {
			fmt.Println(link)
		}
		fmt.Println("\n\n")

		fmt.Println("Photos found: ")
		for i, photo := range photos {
            downloadURL := baseURL + photo

            fmt.Println(downloadURL)

            // Download photo
            resp, err := http.Get(downloadURL)
            if err != nil {
                log.Fatal(err)
            }

            defer resp.Body.Close()

            // Create a directory
            os.MkdirAll("karina", os.ModePerm)

            // Create the file

            out, err := os.Create(fmt.Sprintf("karina/photo%d.jpg", i))
            if err != nil {
                log.Fatal(err)
            }

            defer out.Close()

            // Write the body to file
            _, err = io.Copy(out, resp.Body)
            if err != nil {
                log.Fatal(err)
            }

            fmt.Println("Photo downloaded: ", photo)

		}
	})

	c.Visit(baseURL + "/kpics")

	fmt.Println("Hello, World!")
}
