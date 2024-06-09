package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func test(c echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, "Hello, "+name)
}

func main() {
	e := echo.New()

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	// Format like: INFO "GET / HTTP/1.1" 200 13
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogMethod:       true,
		LogProtocol:     true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("method", v.Method).
				Str("URI", v.URI).
				Str("protocol", v.Protocol).
				Int("status", v.Status).
                Int64("response_size", v.ResponseSize).
				Msg("request")

			return nil
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<form action='/test' method='post'><input type='text' name='name' /><input type='submit' /></form>")
	})

	e.POST("/test", func(c echo.Context) error {
		name := c.FormValue("name")

		return c.String(http.StatusOK, "Hello, "+name)
	})

	e.GET("/test", test)

	e.Logger.Fatal(e.Start(":1323"))
}
