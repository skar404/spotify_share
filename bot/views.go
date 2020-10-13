package bot

import (
	"fmt"

	"github.com/labstack/gommon/log"

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

func CommandHandler(update *telegram.Update) {
	log.Infof("send message: %s", update.Message.Text)

	command, err := getCommand(update.Message.Text)
	// Обрабатываем только команды, если нет то скипаем
	if err != nil {
		return
	}
	// TODO добавить получиния пользователя из базы
	user := FakeUser{
		Uuid:         "MongoDB uuid",
		TelegramId:   1234,
		Token:        "spotify token",
		RefreshToken: "spotify refresh token",
		Active:       true,
		Block:        false,
	}

	if user.Block == true {
		return
	}

	if user.Active == false {
		// TODO go to start
		return
	}

	m := ""
	switch command.Name {
	case "start":
		m = TemplateMessageStart
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

	_ = telegram.TgClient.SendMessage(update.Message.Chat.Id, m)
}
