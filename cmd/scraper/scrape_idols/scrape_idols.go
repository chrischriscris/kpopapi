package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/gocolly/colly"
)

const baseURL = "https://kpopping.com"

func getGroupsFromString(groups string) []string {
	groups = groups[2 : len(groups)-2]
	return strings.Split(groups, ", ")
}

func scrapeIdols(url string) map[string][]string {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting scraping: ", r.URL)
	})

	// idols[idol_name] = [group1, group2, ...]
	idols := make(map[string][]string)
	c.OnHTML("div.item", func(e *colly.HTMLElement) {
		idol := e.ChildText("a")

		// Groups in the format ( <group1>, <group2>, ... )
		groups := e.ChildText("span")
		groupsArr := make([]string, 0)

		if groups != "" {
			groupsArr = getGroupsFromString(groups)
		}

		idols[idol] = groupsArr
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Scraped ", len(idols), " idols")
	})

	c.Visit(url)

	return idols
}

func createIdolHelper(
	q *repository.Queries,
	ctx context.Context,
	idolName string,
	gender string,
) (int32, error) {
	idol, err := q.CreateIdolWithGroupMinimal(ctx, repository.CreateIdolWithGroupMinimalParams{
		StageName: idolName,
		Gender:    gender,
	})
	if err != nil {
		return 0, err
	}

	return idol.ID, err
}

func createGroupHelper(
	qtx *repository.Queries,
	ctx context.Context,
	groupName string,
	groupType string,
) (int32, error) {
	group, err := qtx.CreateGroupMinimal(ctx, repository.CreateGroupMinimalParams{
		Name: groupName,
		Type: groupType,
	})
    if err != nil {
        return 0, err
    }

	return group.ID, err
}

func addIdolToGroup(
	qtx *repository.Queries,
	ctx context.Context,
	groupName string,
	idolID int32,
) error {
    groupName = strings.TrimSpace(groupName)
    groupID := int32(0)
	group, err := qtx.GetGroupByName(ctx, groupName)
	if err != nil {
        groupID, err = createGroupHelper(qtx, ctx, groupName, "UN")
        if err != nil {
            return err
        }
	} else {
        groupID = group.ID
    }

	_, err = qtx.AddMemberToGroup(ctx, repository.AddMemberToGroupParams{
		GroupID: groupID,
		IdolID:  idolID,
	})

	return err
}

func loadToDB(groups map[string]map[string][]string) {
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
	for gender, idols := range groups {
		for idolName, idolGroups := range idols {
			idolID, err := createIdolHelper(qtx, ctx, idolName, gender)
			if err != nil {
				log.Fatalf("Unable to create idol: %v\n", err)
			}

			for _, group := range idolGroups {
				err := addIdolToGroup(qtx, ctx, group, idolID)
				if err != nil {
					log.Fatalf("Unable to add idol %s to group %s: %v\n", idolName, group, err)
				}
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
	}
}

func main() {
	baseIdolsURL := baseURL + "/profiles/the-idols"

	idols := make(map[string]map[string][]string)
	idols["F"] = scrapeIdols(baseIdolsURL + "/women")
	idols["M"] = scrapeIdols(baseIdolsURL + "/men")

	loadToDB(idols)
}
