package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/chrischriscris/kpopapi/internal/scraper"

	"github.com/gocolly/colly"
)

func scrapeGroups(url string) []string {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting scraping: ", r.URL)
	})

	groups := make([]string, 0)
	c.OnHTML("div.item", func(e *colly.HTMLElement) {
		group := e.ChildText("a")
		groups = append(groups, group)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Scraped ", len(groups), " groups")
	})

	c.Visit(url)

	return groups
}

func loadToDB(groups map[string][]string) {
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

	for groupType, groupNames := range groups {
		for _, groupName := range groupNames {
			_, err := qtx.CreateGroupMinimal(ctx, repository.CreateGroupMinimalParams{
				Name: groupName,
				Type: groupType,
			})
			if err != nil {
				fmt.Println("Unable to insert group: ", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
	}
}

func main() {
	baseGroupsURL := scraperutils.BaseURL + "/profiles/the-groups"

	groups := make(map[string][]string)
	groups["GG"] = scrapeGroups(baseGroupsURL + "/women")
	groups["BG"] = scrapeGroups(baseGroupsURL + "/men")
	groups["CE"] = scrapeGroups(baseGroupsURL + "/coed")

	loadToDB(groups)
}
