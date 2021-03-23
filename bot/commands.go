package bot

import (
	"github.com/labstack/gommon/log"
	"github.com/skar404/spotify_share/libs"
	"github.com/skar404/spotify_share/model"
	"github.com/skar404/spotify_share/spotify"
	"github.com/skar404/spotify_share/telegram"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandContext struct {
	update  *telegram.Update
	user    *model.User
	command *Command

	DB *mongo.Database
}

func (c *CommandContext) StartCommand() {
	var rm telegram.InlineKeyboardReq

	m := TemplateMessageStart
	url := spotify.OAuthClient.GetOAthUrl(libs.CreateJWT(c.user.Telegram.Id))
	rm.InlineKeyboard = [][]telegram.InlineKeyboardButtonReq{{{
		Text: "Войти через Spotify",
		Url:  url,
	}}}
	rm.Ready()

	// FIXME попарвить когда идею с ARGS
	//if len(c.command.Args) > 0 {
	//	m = fmt.Sprintf(TemplateMessageStart, "@spotify_share_bot")
	//}
	_ = telegram.Client.SendMessage(c.update.Message.Chat.Id, m, &rm)
}

func (c *CommandContext) Help() {
	text := `Этот бот позволяет делиться audio пользователям spotify
для этого нужно ввести имя боа @spotify\_share\_bot 
и выбрать трек.

Список треков получается из вашего Spotify  

Если у вас есть вопросы, идеи, хотите помочь или нашли баг/опечатку, то напишите: 
 - автору @SaladMen
 - задачи на github https://github.com/skar404/spotify\_share/discussions/11

Новости бота @spotify\_share`

	if err := telegram.Client.SendMessage(c.update.Message.Chat.Id, text, nil); err != nil {
		log.Errorf("error send message err=%s", err)
	}

}
