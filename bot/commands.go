package bot

import (
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

}
