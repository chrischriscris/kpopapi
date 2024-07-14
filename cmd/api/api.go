package main

import (
	"html/template"
	"io"
	"os"

	index "github.com/chrischriscris/kpopapi/internal/handlers"
	"github.com/chrischriscris/kpopapi/internal/scheduler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
    _ "github.com/joho/godotenv/autoload"
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

func main() {
	e := echo.New()
    e.Static("/static", "public/static")

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

	e.GET("/", index.Index)
	e.GET("/idols/random", index.Random)
	e.GET("/idols", index.Idol)
    e.GET("/health", index.Health)

	e.POST("/fetch-new-images", index.FetchNewImages)

    s := scheduler.KPopApiScheduler()
    if os.Getenv("APP_ENV") == "dev" {
        s.Disable()
    }
    s.Start()
    defer s.Shutdown()

	e.Logger.Fatal(e.Start(":8080"))
}
