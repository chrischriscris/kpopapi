package helpers

import (
	"context"
	"log"
	"os"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (context.Context, *pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return ctx, conn, err
}

func BeginTransaction(ctx context.Context, conn *pgx.Conn) (pgx.Tx, *repository.Queries, error) {
	instance := repository.New(conn)
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("Unable to start transaction: %v\n", err)
	}

	qtx := instance.WithTx(tx)
    return tx, qtx, err
}
