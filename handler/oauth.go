package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/skar404/spotify_share/libs"
	"github.com/skar404/spotify_share/model"
	"github.com/skar404/spotify_share/spotify"

	//"github.com/skar404/spotify_share/spotify"
	"github.com/skar404/spotify_share/telegram"
)

// OAuthSpotify method login spotify and redirect ot bot
func (h *Handler) OAuthSpotify(c echo.Context) (err error) {
	// Redirect Spotify to app
	// - validate oauth
	// - get token
	// - redirect to bot

	botUrl := "https://t.me/spotify_share_bot"
	conn := model.Conn{
		DB: h.DB,
	}

	code := c.FormValue("code")
	// TODO may be start code in JWT and set user id ...
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

	user, err := conn.GetUser(userInfo.UserId)
	if err != nil {
		// may be return return to localhost or
		// logging error
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:USER")
	}

	token, err := spotify.OAuthClient.GetAccessOrRefreshToken(code)
	if err != nil || token.RefreshToken == "" {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:TOKEN")
	}

	// Save token in OldDB
	spotifyToken := model.Spotify{Token: model.SpotifyToken{
		Refresh: token.RefreshToken,
		User:    token.TokenReq.AccessToken,
	}}

	err = conn.UpdateSpotifyToken(&user.Id, &spotifyToken)
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, botUrl+"?start=ERROR:UPDATE")
	}

	err = telegram.TgClient.SendMessage(userInfo.UserId, "вы успешно вошли", nil)
	return c.Redirect(http.StatusMovedPermanently, botUrl)
}

func (h *Handler) Ping(c echo.Context) (err error) {
	ctx := context.Background()
	if err := h.DBMongoDB.Ping(ctx, readpref.Primary()); err != nil {
		errJson := c.JSON(http.StatusBadGateway, map[string]bool{"ok": false})
		if errJson != nil {
			log.Error(errJson)
		}
		return err
	}

	collUser := h.DB.Collection("user")

	user := model.User{
		Telegram: model.Telegram{
			Id:    123,
			Login: "app",
		},
	}

	result, err := collUser.InsertOne(ctx, &user)
	if err != nil {
		log.Error("error create user err=", err)
		_ = c.JSON(http.StatusBadGateway, map[string]bool{"ok": false})
		return nil
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.Id = oid
	} else {
		_ = c.JSON(http.StatusBadGateway, map[string]bool{"ok": false})
		return nil
	}

	user.Active = true
	user.UpdateAt = time.Now()
	user.Spotify.Token.Refresh = "radsa a12 "

	log.Info(collUser.UpdateOne(ctx,
		bson.M{"_id": &user.Id},
		bson.D{
			{"$set", &user},
		},
	))

	res := collUser.FindOne(ctx, bson.M{"telegram.id": 12223})

	findUser := model.User{}

	err = res.Decode(&findUser)

	_is := errors.Is(err, mongo.ErrNoDocuments)
	_ = _is

	return c.JSON(http.StatusOK, map[string]interface{}{"ok": true})
}
