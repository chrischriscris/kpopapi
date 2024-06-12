package main

import (
	"html/template"
	"io"
	"net/http"
	"os"

	dbutils "github.com/chrischriscris/kpopapi/internal/db"
	"github.com/chrischriscris/kpopapi/internal/db/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

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

func main() {
	e := echo.New()

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	e.Renderer = NewTemplate()

	// Example format: 2:58PM INF "GET / HTTP/1.1" 200 13
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogMethod:       true,
		LogProtocol:     true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().Msgf("\"%s %s %s %d\" %d", v.Method, v.URI, v.Protocol, v.Status, v.ResponseSize)

			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.GET("/", Index)
	e.GET("/random", Random)
	e.GET("/idols", Idol)

	e.Logger.Fatal(e.Start(":8080"))
}
