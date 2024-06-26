package dbutils

import (
	"context"
	"log"
	"os"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
)

func ConnectDB() (context.Context, *pgx.Conn, error) {
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

