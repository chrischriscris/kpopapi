package main

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"reflect"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgresql://postgres:password@localhost:5432/kpop")

	queries := repository.New(db)

	idols, _ := queries.ListIdols(ctx)
	log.Println("List of idols")
	log.Println(idols)

	insertedIdol, err := queries.CreateIdol(ctx, repository.CreateIdolParams{
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
