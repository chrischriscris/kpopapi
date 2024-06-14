package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
}

type service struct {
	db *pgxpool.Pool
}

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := pgxpool.NewWithConfig(context.Background(), Config())
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{db }
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	// Database is up, add more stats
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stat()
	stats["open_connections"] = strconv.Itoa(int(dbStats.AcquiredConns()))
	stats["idle"] = strconv.Itoa(int(dbStats.IdleConns()))

	// TODO: Improve this with more stats and handle stress cases
	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() {
	log.Printf("Disconnecting from database: %s", database)
	s.db.Close()
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
