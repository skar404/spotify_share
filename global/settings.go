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
	RedirectUri  = getEnv("REDIRECT_URI", "http://localhost/spotify")

	AppMode = os.Getenv("APP_MOD")

	JWTToken = []byte(getEnv("JWT_TOKEN", "TEST_TOKEN"))
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
