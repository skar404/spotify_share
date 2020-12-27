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
	initStopSignal()
	// App env
	webhookToken := global.WebhookToken
	telegramToken := global.TelegramToken

	clientId := global.ClientId
	clientSecret := global.ClientSecret

	appMode := global.AppMode

	// Database connection
	url := fmt.Sprintf("mongodb://%s:%s@%s/%s?replicaSet=%s&tls=true&tlsCaFile=%s",
		global.DBUser,
		global.DBPass,
		global.DBHost,
		global.DBName,
		global.DBRs,
		global.DBCACERT)

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	defer conn.Disconnect(context.Background())

	db := conn.Database(global.DBName)
	log.Info("conn to OldDB", db.Name())

	// Initialize handler
	h := &handler.Handler{
		DBMongoDB: conn,
		DB:        db,
	}

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
			runGetUpdate(telegramToken, h)
			lockChanel <- true
		}()
	}

	if appMode != "CLI" {
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
		_ = fmt.Errorf("Error create token")
		return
	}

	fmt.Println("token", token, refreshToken)
}

func runGetUpdate(telegramToken string, h *handler.Handler) {
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

			bot.BotRouter(&item, h)
			updateId = item.UpdateId + 1
		}
	}
}

func runHttpServer(webhookToken string, handler *handler.Handler) {
	e := echo.New()
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

	e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}
