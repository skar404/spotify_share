package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/skar404/spotify_share/bot"
	"github.com/skar404/spotify_share/commands"
	"github.com/skar404/spotify_share/telegram"
)

func main() {
	// App env
	webhookToken := os.Getenv("TELEGRAM_WEBHOOK_TOKEN")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	appMode := os.Getenv("APP_MOD")

	if appMode == "CLI" {
		runCLI(clientId, clientSecret)
	} else if appMode == "WEB" {
		runHttpServer(webhookToken)
	} else if appMode == "GET_UPDATES" {
		runGetUpdate(telegramToken)
	}
}

func runCLI(clientId, clientSecret string) {
	token, refreshToken, err := commands.CreateToken(clientId, clientSecret)
	if err != nil {
		_ = fmt.Errorf("Error create token")
		return
	}

	fmt.Println("token", token, refreshToken)
}

func runGetUpdate(telegramToken string) {
	tg, _ := telegram.Init(telegramToken)

	updateId := 0
	for true {
		raw, err := tg.GetUpdates(updateId)

		if err != nil {
			log.Println(err)
			continue
		}

		for _, item := range raw.Result {
			if updateId > item.UpdateId {
				continue
			}

			fmt.Println(fmt.Sprintf("send message: %s", item.Message.Text))
			updateId = item.UpdateId + 1
		}
	}
}

func runHttpServer(webhookToken string) {
	e := echo.New()

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET(
		fmt.Sprintf("/api/telegram/webhook/%s", webhookToken),
		bot.Router,
	)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
