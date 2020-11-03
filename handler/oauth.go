package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/skar404/spotify_share/libs"
	"github.com/skar404/spotify_share/spotify"
	"github.com/skar404/spotify_share/telegram"
)

func (h *Handler) OAuthSpotify(c echo.Context) (err error) {
	// Redirect Spotify to app
	// - validate oauth
	// - get token
	// - redirect to bot

	botUrl := "https://t.me/spotify_share_bot"

	code := c.FormValue("code")
	if code == "" {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:CODE")
	}

	state := c.FormValue("state")
	if state == "" {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:STATE")
	}

	userInfo, err := libs.DecodeJWT(state)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:JWT")
	}

	token, err := spotify.OAuthClient.GetAccessOrRefreshToken(code)
	if err != nil || token.RefreshToken == "" {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:TOKEN")
	}

	// Save token in DB
	_ = token
	_ = userInfo
	err = telegram.TgClient.SendMessage(userInfo.UserId, "вы успешно вошли", nil)

	return c.Redirect(http.StatusMovedPermanently, "https://t.me/spotify_share_bot")
}
