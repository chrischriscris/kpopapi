package index

import (
	"fmt"
	"net/http"
	"os"

	dbutils "github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	images "github.com/chrischriscris/kpopapi/internal/scraper/kpopping"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

type IndexData struct {
	Image string
	Idols []repository.Idol
    NumberOfImages int64
}

func NewIndexData(image string, idols []repository.Idol, n int64) IndexData {
	return IndexData{
		Image: image,
		Idols: idols,
        NumberOfImages: n,
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

    n, err := queries.GetNumberOfImages(ctx)
    if err != nil {
        n = 0
    } else {
        n -= 1
    }

	return c.Render(http.StatusOK, "index", NewIndexData(image.Url, nil, n))
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
	name := c.QueryParam("name")
    if name == "" {
        return c.Render(http.StatusOK, "idol", nil)
    }

	ctx, conn, err := dbutils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := repository.New(conn)
	idols, err := queries.GetIdolsByNameLike(ctx, pgtype.Text{String: name, Valid: true})
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "idols", idols)
}

// This can be better, not loading the .env file every time
func isAdmin(c *echo.Context) error {
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
