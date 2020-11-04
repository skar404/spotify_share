package bot

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/skar404/spotify_share/handler"
	"github.com/skar404/spotify_share/model"
	"github.com/skar404/spotify_share/telegram"
)

var WarningMessageNotCommand = errors.New("message not command")
var NotValidUser = errors.New("user not found in telegram")

func isCommand(m string) bool {
	if len(m) > 0 && m[0] == '/' && len(m[1:]) > 0 {
		return true
	}
	return false
}

func getCommand(m string) (*Command, error) {
	c := Command{}
	if !isCommand(m) {
		return nil, WarningMessageNotCommand
	}

	mSplit := strings.Split(m, " ")
	c.Name = mSplit[0][1:]
	c.Args = mSplit[1:]

	return &c, nil
}

func GetOrCreateUser(tgUser *telegram.User, handler *handler.Handler) (*model.User, error) {
	if tgUser.Id == 0 {
		return nil, NotValidUser
	}

	conn := model.Conn{
		DB: handler.DB,
	}
	user, err := conn.GetUser(tgUser.Id)
	if err == nil {
		return user, nil
	} else if err != mgo.ErrNotFound {
		return nil, err
	}

	u := &model.User{
		Id: bson.NewObjectId(),

		Telegram: model.Telegram{
			Id:    tgUser.Id,
			Login: tgUser.Username,
		},
		Spotify:  model.Spotify{},
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Active:   true,
	}
	u, err = conn.CreateUser(u)

	if err != nil {
		return nil, err
	}
	return u, nil
}
