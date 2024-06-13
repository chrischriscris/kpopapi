package index

import (
	"net/http"

	dbutils "github.com/chrischriscris/kpopapi/internal/db"
	images "github.com/chrischriscris/kpopapi/internal/scraper/kpopping"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5/pgtype"
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

func FetchNewImages(c echo.Context) error {
    images.ScrapeImages()

    return c.String(http.StatusOK, "Successfully fetched new images")
}
