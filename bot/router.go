package bot

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Router(c echo.Context) error {
	// Command:
	// -> start
	// -> login
	// -> logout
	// -> bot command hook

	return c.String(http.StatusOK, "")
}
