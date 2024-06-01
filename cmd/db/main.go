package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"

	"github.com/chrischriscris/kpopapi/internal/db/tutorial"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	queries := tutorial.New(db)

	idols, _ := queries.ListIdols(ctx)
	log.Println("List of idols")
	log.Println(idols)

	insertedIdol, err := queries.CreateIdol(ctx, tutorial.CreateIdolParams{
		StageName:  "Chaeyoung",
		Name:       "Son Chae-young",
		Gender:     "F",
		IdolInfoID: sql.NullInt64{Int64: 1, Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted idolID: ", insertedIdol.ID)

	fetchedIdol, err := queries.GetIdol(ctx, insertedIdol.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Is insertedIdol equal to fetchedIdol?")
	log.Println(reflect.DeepEqual(insertedIdol, fetchedIdol))

	log.Println("What is fetchedIdol?")
	log.Println(fetchedIdol)
}
