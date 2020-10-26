package main

import (
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/skar404/spotify_share/bot"
	"github.com/skar404/spotify_share/commands"
	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/telegram"
)

func main() {
	initStopSignal()
	// App env
	webhookToken := global.WebhookToken
	telegramToken := global.TelegramToken

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	appMode := global.AppMode

	InitGlobal(telegramToken)

	lockChanel := make(chan bool)
	if appMode == "CLI" {
		runCLI(clientId, clientSecret)
	} else if appMode == "WEB" {
		// set webhook
	} else if appMode == "GET_UPDATES" {
		// only dev mode:
		// ... run `for true` and lock web server
		go func() {
			log.Info("create goroutines")
			runGetUpdate(telegramToken)
			lockChanel <- true
		}()
	}
	runHttpServer(webhookToken)

	<-lockChanel
}

func InitGlobal(telegramToken string) {
	telegram.TgClient, _ = telegram.Init(telegramToken)

}

func initStopSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Infof("stop apps sig=%v", sig)
			os.Exit(1)
		}
	}()
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
	tg := &telegram.TgClient

	updateId := 0
	for true {
		raw, err := tg.GetUpdates(updateId)

		if err != nil {
			stdLog.Println(err)
			continue
		}

		for _, item := range raw.Result {
			if updateId > item.UpdateId {
				continue
			}

			bot.CommandHandler(&item)
			updateId = item.UpdateId + 1
		}
	}
}

func OAuthSpotify(c echo.Context) error {
	// Redirect Spotify to app
	// - validate oauth
	// - get token
	// - redirect to bot
	return c.Redirect(http.StatusMovedPermanently, "https://t.me/spotify_share_bot?start=token")
}

func runHttpServer(webhookToken string) {
	e := echo.New()
	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	// Middleware
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET(
		fmt.Sprintf("/api/telegram/webhook/%s", webhookToken),
		bot.Router,
	)
	e.GET(
		"/oauth",
		OAuthSpotify,
	)

	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
