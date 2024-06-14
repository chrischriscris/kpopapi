package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

const defaultMaxConns = 4
const defaultMinConns = 0
const defaultMaxConnLifetime = time.Hour
const defaultMaxConnIdleTime = time.Minute * 30
const defaultHealthCheckPeriod = time.Minute
const defaultConnectTimeout = time.Second * 5

func Config() *pgxpool.Config {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Acquiring connection pool to the database")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("Successfully released connection pool to the database")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Successfully closed connection pool to the database")
	}

	return dbConfig
}
