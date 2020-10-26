package global

import (
	"os"
)

var (
	WebhookToken  = os.Getenv("TELEGRAM_WEBHOOK_TOKEN")
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")

	AppSpotifyScope = []string{
		"user-read-recently-played",
		"user-read-currently-playing",
		"app-remote-control",
		"streaming"}
	ClientId     = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")

	AppMode = os.Getenv("APP_MOD")
)
