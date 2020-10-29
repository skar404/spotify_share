package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) OAuthSpotify(c echo.Context) (err error) {
	// Redirect Spotify to app
	// - validate oauth
	// - get token
	// - redirect to bot
	return c.Redirect(http.StatusMovedPermanently, "https://t.me/spotify_share_bot?start=token")
}
