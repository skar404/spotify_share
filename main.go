package main

import (
	"context"
	"fmt"
	stdLog "log"
	"os"
	"os/signal"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/skar404/spotify_share/bot"
	"github.com/skar404/spotify_share/commands"
	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/handler"
	"github.com/skar404/spotify_share/telegram"
)

func main() {
	mode := global.AppMode
	log.Infof("starting app, config: mode=%s, debug=%t", mode, global.Debug)

	initStopSignal()
	// App env
	webhookToken := global.WebhookToken
	telegramToken := global.TelegramToken

	clientId := global.ClientId
	clientSecret := global.ClientSecret

	url := constructDBUrl()
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	defer conn.Disconnect(context.Background())
	db := conn.Database(global.DBName)

	// Initialize handler
	h := &handler.Handler{
		DBConn: conn,
		DB:     db,
	}

	lockChanel := make(chan bool)
	if mode == "CLI" {
		runCLI(clientId, clientSecret)
	} else if mode == "WEB" {
		// set webhook
	} else if mode == "GET_UPDATES" {
		// only dev mode:
		// ... run `for true` and lock web server
		go func() {
			runGetUpdate(telegramToken, h)
			lockChanel <- true
		}()
	}

	if mode != "CLI" {
		runHttpServer(webhookToken, h)
		<-lockChanel
	}
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
		fmt.Printf("error create token %e\n", err)
		return
	}

	fmt.Printf("token=%s, refreshTo¬ken=%s\n", token, refreshToken)
}

func runGetUpdate(telegramToken string, h *handler.Handler) {
	tg := &telegram.Client

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

			bot.BotRouter(&item, h)
			updateId = item.UpdateId + 1
		}
	}
}

func constructDBUrl() string {
	url := fmt.Sprintf("mongodb://%s:%s@%s/%s?",
		global.DBUser,
		global.DBPass,
		global.DBHost,
		global.DBName)

	if global.DBAuthSource != "" {
		url = fmt.Sprintf("%s&authSource=%s", url, global.DBAuthSource)
	}

	if global.DBRs != "" {
		url = fmt.Sprintf("%s&replicaSet=%s", url, global.DBRs)
	}

	if global.DBCACERT != "" {
		url = fmt.Sprintf("%s&tls=true&tlsCaFile=%s", url, global.DBCACERT)
	}
	return url
}

func runHttpServer(webhookToken string, handler *handler.Handler) {
	e := echo.New()
	e.Debug = global.Debug

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	// Middleware
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/spotify", handler.OAuthSpotify)
	e.GET("/ping", handler.Ping)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:1323", global.AppHost)))
}
