package admin

import (
	"fmt"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
    fmt.Println("admin")
	return c.Render(http.StatusOK, "admin", nil)
}
