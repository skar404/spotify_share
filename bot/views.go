package bot

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/skar404/spotify_share/handler"
	_ "github.com/skar404/spotify_share/handler"
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

func BotRouter(update *telegram.Update, handler *handler.Handler) {
	if update.Message.MessageId != 0 {
		CommandHandler(update, handler)
	} else if update.InlineQuery.Id != "" {
		InlineQueryHandler(update, handler)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func InlineQueryHandler(update *telegram.Update, handler *handler.Handler) {
	user, err := GetOrCreateUser(&update.InlineQuery.From, handler)
	if err != nil {
		return
	}
	var tmpList []interface{}

	token, _ := spotify.OAuthClient.RefreshToken(user.Spotify.Token.Refresh)

	api := spotify.ApiClient.SetUserToken(token.AccessToken)
	r, _ := api.GetHistory()
	playNow, err := api.GetPlayNow()

	if err == nil {
		tmpList = append(tmpList, map[string]interface{}{
			"type":  "photo",
			"id":    fmt.Sprintf("%v %v", time.Now().Unix(), RandStringBytes(10)),
			"title": playNow.Item.Name,
			"description": fmt.Sprintf("%s",
				playNow.Item.Artists[0].Name),
			"is_personal": true,
			//"input_message_content": map[string]interface{}{
			//	"message_text": fmt.Sprintf("test ![img](%s)", playNow.Item.Album.Images[len(playNow.Item.Album.Images)-1].URL),
			//	"parse_mode":   "Markdown",
			//},
			"caption": fmt.Sprintf("Name: ***%s***\nArtist: ***%s***\nAlbum: ***%s***\ndebug info: inline ID=%s",
				playNow.Item.Name,
				playNow.Item.Artists[0].Name,
				playNow.Item.Album.Name,
				update.InlineQuery.Id),
			"parse_mode": "Markdown",
			"photo_url":  playNow.Item.Album.Images[0].URL,
			"reply_markup": map[string]interface{}{
				"inline_keyboard": [][]map[string]interface{}{{{
					"text":          "Play in Spotify",
					"callback_data": "https://google.com",
				}}},
			},
			"thumb_url": playNow.Item.Album.Images[len(playNow.Item.Album.Images)-1].URL,
		})
	}

	for _, value := range r.Items {
		tmpList = append(tmpList, map[string]interface{}{
			"type":        "photo",
			"id":          fmt.Sprintf("%v %v", time.Now().Unix(), RandStringBytes(10)),
			"title":       value.Track.Name,
			"description": fmt.Sprintf("%s"),
			"caption": fmt.Sprintf("%s\n%s\nInline ID=%s",
				value.Track.Artists[0].Name,
				value.Track.Album.Name,
				update.InlineQuery.Id),
			"is_personal": true,
			"photo_url":   value.Track.Album.Images[0].URL,
			//"input_message_content": map[string]interface{}{
			//	"message_text": "test",
			//	"parse_mode":   "Markdown",
			"parse_mode": "Markdown",
			"reply_markup": map[string]interface{}{
				"inline_keyboard": [][]map[string]interface{}{{{
					"text":          "Play in Spotify",
					"callback_data": "https://google.com",
				}}},
			},
			"thumb_url": value.Track.Album.Images[len(value.Track.Album.Images)-1].URL,
		})
	}

	_ = telegram.TgClient.AnswerInlineQuery(update.InlineQuery.Id, tmpList)
	_ = user
	_ = r
	_ = token
}

func CommandHandler(update *telegram.Update, handler *handler.Handler) {
	user, err := GetOrCreateUser(&update.Message.From, handler)
	if err != nil {
		return
	}

	command, err := getCommand(update.Message.Text)
	// Обрабатываем только команды, если нет то скипаем
	if err != nil {
		return
	}

	if user.Active == false {
		// skip block user ...
		return
	}
	log.Infof("send message: %s", update.Message.Text)

	commands := CommandContext{
		update:  update,
		user:    user,
		command: command,
	}

	switch command.Name {
	case "start":
		commands.StartCommand()
	case "help":

	case "logout":

	default:

	}
}
