package index

import (
	"fmt"
	"log"
	"net/http"
	"os"

	dbutils "github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	images "github.com/chrischriscris/kpopapi/internal/scraper/kpopping"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type IndexData struct {
	Image string
	Idols []repository.Idol
}

func NewIndexData(image string, idols []repository.Idol) IndexData {
	return IndexData{
		Image: image,
		Idols: idols,
	}
}

func Index(c echo.Context) error {
	ctx, conn, err := dbutils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)
	image, err := queries.GetRandomImage(ctx)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index", NewIndexData(image.Url, nil))
}

func Random(c echo.Context) error {
	ctx, conn, err := dbutils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)
	image, err := queries.GetRandomImage(ctx)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "image", image.Url)
}

func Idol(c echo.Context) error {
	ctx, conn, err := dbutils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	name := c.QueryParam("name")
	queries := repository.New(conn)
	idols, err := queries.GetIdolsByNameLike(ctx, pgtype.Text{String: name, Valid: true})
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "idols", idols)
}

// This can be better, not loading the .env file every time
func isAdmin(c *echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

    auth := (*c).Request().Header.Get("Authorization")
    if auth != os.Getenv("SECRET") {
        return fmt.Errorf("Unauthorized")
    }

    return nil
}

func FetchNewImages(c echo.Context) error {
    err := isAdmin(&c)
    if err != nil {
        return c.String(http.StatusUnauthorized, "Unauthorized")
    }

    n := images.ScrapeImages()
    msg := fmt.Sprintf("Successfully fetched %d new images", n)
    return c.String(http.StatusOK, msg)
}

func Health(c echo.Context) error {
    return c.String(http.StatusOK, "OK")
}
