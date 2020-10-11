package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func login(c echo.Context) error {
	// -> start
	// -> login
	// -> logout
	// -> bot command hook
	return c.String(http.StatusOK, "Hello, World!")
}

func auth() {

}

func main() {
	// App env
	webhookToken := os.Getenv("TELEGRAM_WEBHOOK_TOKEN")
	//clientId := os.Getenv("CLIENT_ID")
	//clientSecret := os.Getenv("CLIENT_SECRET")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET(
		fmt.Sprintf("/api/telegram/webhook/%s", webhookToken),
		login,
	)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
