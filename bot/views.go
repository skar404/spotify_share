package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/skar404/spotify_share/handler"
	_ "github.com/skar404/spotify_share/handler"
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

func BotRouter(update *telegram.Update, handler *handler.Handler) {

	// FIXME mey be create struct
	if update.Message.MessageId != 0 {
		CommandHandler(update, handler)
	} else if update.InlineQuery.Id != "" {
		InlineQueryHandler(update, handler)
	} else if update.CallbackQuery.Id != "" {
		CallbackQueryHandler(update, handler)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CallbackQueryHandler(update *telegram.Update, handler *handler.Handler) {
	callback := update.CallbackQuery
	conn := model.Conn{
		DB: handler.DB,
	}
	user, err := conn.GetUser(callback.From.Id)

	data := telegram.AnswerCallbackReq{}
	if err != nil || user.Spotify.Token == nil {
		data.Url = "t.me/spotify_share_bot?start=LOGIN"
	}

	_ = telegram.Client.AnswerCallbackQuery(callback.Id, &data)

	if user == nil || user.Spotify.Token == nil {
		return
	}
	token, _ := spotify.OAuthClient.RefreshToken(user.Spotify.Token.Refresh)

	api := spotify.ApiClient.SetUserToken(token.AccessToken)

	splitData := strings.SplitN(callback.Data, "::", 2)

	if len(splitData) != 2 {
		log.Info("skip ist=", splitData, ", len=", len(splitData))
		return
	}
	log.Infof("method %s", splitData[0])

	if splitData[0] == "PLAY" {
		// Пока не придумал как можно выклюить трек и сохранить контекст который до этого слушал пользоватлеь...
		// идея:
		// # возможно стоит потестить api очереди и дествовать по алгоритму:
		// - получать контекст плеира
		// - включать трек
		// - подмешивать контекст (желательно не тераю очередь)

		// song, err := api.GetPlayNow()
		// context := api.AddQueue(song.Item.URI)
		if err := api.Play(splitData[1]); err != nil {
			log.Infof("play error=%s", err)
		}

		//_ = api.AddQueue(splitData[1])
	} else if splitData[0] == "ADD" {
		_ = api.AddQueue(splitData[1])
	}
}

// TODO move to libs
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// TODO refactoring !!!
func InlineQueryHandler(update *telegram.Update, handler *handler.Handler) {
	user, err := GetOrCreateUser(&update.InlineQuery.From, handler)

	if err != nil {
		return
	}
	if user.Spotify.Token != nil {
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
					"inline_keyboard": [][]map[string]interface{}{{
						{
							"text":          "Play",
							"callback_data": fmt.Sprintf("PLAY::%s", playNow.Item.URI),
						},
						{
							"text":          "Add",
							"callback_data": fmt.Sprintf("ADD::%s", playNow.Item.URI),
						},
						//{
						//	"text":          "Sync",
						//	"callback_data": fmt.Sprintf("SYNC:%s", playNow.Item.URI),
						//},
					}},
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
					"inline_keyboard": [][]map[string]interface{}{{
						{
							"text":          "Play",
							"callback_data": fmt.Sprintf("PLAY::%s", value.Track.URI),
						},
						{
							"text":          "Add",
							"callback_data": fmt.Sprintf("ADD::%s", value.Track.URI),
						}}},
				},
				"thumb_url": value.Track.Album.Images[len(value.Track.Album.Images)-1].URL,
			})
		}
		err = telegram.Client.AnswerInlineQuery(update.InlineQuery.Id, tmpList)
		if err != nil {
			log.Error("AnswerInlineQuery err=", err)
		}
	} else {
		r := map[string]interface{}{
			"is_personal":         true,
			"switch_pm_text":      "login in spotify ...",
			"switch_pm_parameter": "inline_redirect",
		}
		err = telegram.Client.AnswerInlineQueryTmp(update.InlineQuery.Id, r)
		log.Info("app")

	}
}

func CommandHandler(update *telegram.Update, handler *handler.Handler) {
	user, err := GetOrCreateUser(&update.Message.From, handler)
	if err != nil {
		log.Error("error create user err=", err)
		return
	}

	command, err := getCommand(update.Message.Text)
	// обрабатываем только команды, если нет то скипаем
	if err != nil {
		log.Info("skip text err=", err)
		return
	}

	if user.Active == false {
		// skip block user ...
		log.Info("skip block user ...=")
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
