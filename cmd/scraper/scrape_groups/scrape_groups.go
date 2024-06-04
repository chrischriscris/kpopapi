package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/gocolly/colly"
)

const baseURL = "https://kpopping.com"

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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	instance := repository.New(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)

	qtx := instance.WithTx(tx)
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
	baseGroupsURL := baseURL + "/profiles/the-groups"

	groups := make(map[string][]string)
	groups["GG"] = scrapeGroups(baseGroupsURL + "/women")
	groups["BG"] = scrapeGroups(baseGroupsURL + "/men")
	groups["CE"] = scrapeGroups(baseGroupsURL + "/coed")

    loadToDB(groups)
}