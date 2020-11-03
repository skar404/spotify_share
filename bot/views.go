package bot

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/skar404/spotify_share/handler"
	_ "github.com/skar404/spotify_share/handler"
	"github.com/skar404/spotify_share/libs"
	"github.com/skar404/spotify_share/model"
	"github.com/skar404/spotify_share/spotify"
	"github.com/skar404/spotify_share/telegram"
)

type FakeUser struct {
	Uuid         string
	TelegramId   int64
	Token        string
	RefreshToken string

	Active bool
	Block  bool
}

func GetOrCreateUser(update *telegram.Update, handler *handler.Handler) (*model.User, error) {
	var err error
	db := handler.DB.Clone()
	defer db.Close()

	u := &model.User{Id: bson.NewObjectId()}

	err = db.DB("bot_db").C("users").Find(bson.M{"telegram.id": update.Message.From.Id}).One(u)
	if err == nil {
		return u, nil
	}

	if err != mgo.ErrNotFound {
		return u, err
	}

	u = &model.User{
		Id: bson.NewObjectId(),

		Telegram: model.Telegram{
			Id:    update.Message.From.Id,
			Login: update.Message.From.Username,
		},
		Spotify:  model.Spotify{},
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Active:   true,
	}
	err = db.DB("bot_db").C("users").Insert(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func CommandHandler(update *telegram.Update, handler *handler.Handler) {

	user, err := GetOrCreateUser(update, handler)
	if err != nil {
		return
	}

	log.Infof("send message: %s", update.Message.Text)

	command, err := getCommand(update.Message.Text)
	// Обрабатываем только команды, если нет то скипаем
	if err != nil {
		return
	}
	// TODO добавить получиния пользователя из базы

	if user.Active == false {
		// TODO go to start
		return
	}

	m := ""
	var rm telegram.InlineKeyboardReq

	switch command.Name {
	case "start":
		m = TemplateMessageStart
		url := spotify.OAuthClient.GetOAthUrl(libs.CreateJWT(user.Telegram.Id))
		rm.InlineKeyboard = [][]telegram.InlineKeyboardButtonReq{{{
			Text: "Войти через Spotify",
			Url:  url,
		}}}
		rm.Ready()

		if len(command.Args) > 0 {
			// TODO валиддация OAuth,
			//    Args[0] - код авторизовался пользователь или нет
			m = fmt.Sprintf(TemplateMessageLogin, "@spotify_share_bot")
		}
	case "help":
		m = TemplateMessageHelp
	case "logout":
		// TODO проверить что токены валиден
		m = TemplateMessageLogout
	default:
		return
	}

	_ = telegram.TgClient.SendMessage(update.Message.Chat.Id, m, &rm)
}
