package main

import (
	"fmt"
	"log"

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
		for _, photo := range photos {
			fmt.Println(photo)
		}
	})

	c.Visit(baseURL + "/kpics")

	fmt.Println("Hello, World!")
}
