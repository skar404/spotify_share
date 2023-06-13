package global

import (
	"os"
	"strconv"
)

var (
	Debug, _ = strconv.ParseBool(getEnv("DEBUG", "false"))
	AppHost  = getEnv("APP_HOST", "0.0.0.0")
	AppMode  = os.Getenv("APP_MOD")

	WebhookToken  = os.Getenv("TELEGRAM_WEBHOOK_TOKEN")
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	BotName       = getEnv("BOT_NAME", "spotify_share_bot")

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

		// Like song
		//"user-library-modify",
	}
	ClientId     = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	RedirectUri  = getEnv("REDIRECT_URI", "http://localhost:1323/spotify")

	DBRs                     = getEnv("DB_RS", "")
	DBName                   = getEnv("DB_NAME", "")
	DBAuthSource             = getEnv("DB_AUTH_SOURCE", "")
	DBHost                   = getEnv("DB_HOST", "")
	DBUser                   = getEnv("DB_USER", "")
	DBPass                   = getEnv("DB_PASS", "")
	DBCACERT                 = getEnv("DB_CACERT", "")
	DBAuthenticationDatabase = getEnv("DBAuthenticationDatabase", "admin")
	JWTToken                 = getEnv("JWT_TOKEN", "TEST_TOKEN")
	JWTTokenByte             = []byte(JWTToken)
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
