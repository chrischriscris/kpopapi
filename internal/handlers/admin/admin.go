package admin

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	dbutils "github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const YYYYMMDD = "2006-01-02"

func getValidGroupTypes() []string {
	return []string{"GG", "BG", "CE"}
}

func Index(c echo.Context) error {
	fmt.Println("admin")
	return c.Render(http.StatusOK, "admin", nil)
}

func addGroupToDB(c echo.Context) error {
	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusOK, "No group name provided")
	}

	groupType := c.FormValue("group-type")
	if !slices.Contains(getValidGroupTypes(), groupType) {
		groupType = "UN"
	}

	debut := c.FormValue("debut")
	parsedDebut, err := time.Parse(YYYYMMDD, debut)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusOK, "Invalid debut date")
	}

	// TODO: Avoid this later by having a connection pool ready
	// Can I do something like Spring @Transactional annotations?
	ctx, conn, err := dbutils.ConnectDB()
	if err != nil {
		return c.String(http.StatusOK, "Unable to connect to database")
	}
	defer conn.Close(ctx)

	tx, qtx, err := dbutils.BeginTransaction(ctx, conn)
	if err != nil {
		return fmt.Errorf("Unable to begin transaction")
	}
	defer tx.Rollback(ctx)

	_, err = qtx.CreateGroupMinimalWithDebut(ctx, repository.CreateGroupMinimalWithDebutParams{
		Name:      name,
		Type:      groupType,
		DebutDate: pgtype.Date{Time: parsedDebut, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusOK, "Error adding group to database")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return c.String(http.StatusOK, "Unable to commit transaction")
	}

	return c.String(http.StatusOK, "Group added successfully!")
}

func addIdolToDB(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented yet")
}

func AddToDB(c echo.Context) error {
	addType := c.FormValue("type")
	if addType == "" {
		return c.String(http.StatusOK, "No type provided")
	}

	switch addType {
	case "group":
		return addGroupToDB(c)
	case "solo":
		return addIdolToDB(c)
	default:
		return c.String(http.StatusOK, "Invalid type")
	}
}
