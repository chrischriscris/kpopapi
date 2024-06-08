package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/chrischriscris/kpopapi/internal/db/helpers"

	"github.com/gocolly/colly"
)

const baseURL = "https://kpopping.com"

type IdolWithGroups struct {
	IdolName string
	Groups   []string
}

func getGroupsFromString(groups string) []string {
	groups = groups[2 : len(groups)-2]
	groupList := strings.Split(groups, ", ")
	for i, group := range groupList {
		groupList[i] = strings.TrimSpace(group)
	}

	return groupList
}

func scrapeIdols(url string) []IdolWithGroups {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting scraping: ", r.URL)
	})

	idols := make([]IdolWithGroups, 0)
	c.OnHTML("div.item", func(e *colly.HTMLElement) {
		idol := e.ChildText("a")

		// Groups in the format ( <group1>, <group2>, ... )
		groups := e.ChildText("span")
		groupsArr := make([]string, 0)

		if groups != "" {
			groupsArr = getGroupsFromString(groups)
		}

		idols = append(idols, IdolWithGroups{
			IdolName: idol,
			Groups:   groupsArr,
		})
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

func loadToDB(groups map[string][]IdolWithGroups) {
	ctx, conn, err := helpers.ConnectDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	tx, qtx, err := helpers.BeginTransaction(ctx, conn)
	if err != nil {
		log.Fatalf("Unable to start transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)

	for gender, idols := range groups {
		for _, el := range idols {
			idolID, err := createIdolHelper(qtx, ctx, el.IdolName, gender)
			if err != nil {
				log.Fatalf("Unable to create idol: %v\n", err)
			}

			for _, group := range el.Groups {
				err := addIdolToGroup(qtx, ctx, group, idolID)
				if err != nil {
					log.Fatalf("Unable to add idol %s to group %s: %v\n", el.IdolName, group, err)
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

	idols := make(map[string][]IdolWithGroups)
	idols["F"] = scrapeIdols(baseIdolsURL + "/women")
	idols["M"] = scrapeIdols(baseIdolsURL + "/men")

	loadToDB(idols)
}
