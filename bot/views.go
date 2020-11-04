package bot

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/skar404/spotify_share/handler"
	_ "github.com/skar404/spotify_share/handler"
	"github.com/skar404/spotify_share/libs"
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
