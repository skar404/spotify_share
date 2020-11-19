package global

import (
	"os"
)

var (
	WebhookToken  = os.Getenv("TELEGRAM_WEBHOOK_TOKEN")
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")

	AppSpotifyScope = []string{
		// Get user history, method: spotify.GetHistory
		"user-read-recently-played",

		// Get user play now, method: GetPlayNow
		"user-read-currently-playing",

		// Get user device,
		// нужно дял выбора активного устройства
		// -- (если музыка сейчас не играет)
		"user-read-playback-state",

		// Play and sync track
		"app-remote-control",
		"streaming",
	}
	ClientId     = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	RedirectUri  = getEnv("REDIRECT_URI", "http://localhost:1323/spotify")

	AppMode = os.Getenv("APP_MOD")

	DBUrl = getEnv("DB_URL", "root:example@localhost")

	JWTToken     = getEnv("JWT_TOKEN", "TEST_TOKEN")
	JWTTokenByte = []byte(JWTToken)
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
