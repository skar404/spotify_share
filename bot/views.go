package bot

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"math/rand"
	"strings"

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

type Bot struct {
	ctx context.Context

	update  *telegram.Update
	handler *handler.Handler
}

func Router(update *telegram.Update, handler *handler.Handler) {
	bot := Bot{
		context.Background(),
		update,
		handler,
	}

	if update.Message.MessageId != 0 {
		bot.CommandHandler()
	} else if update.InlineQuery.Id != "" {
		bot.InlineQueryHandler()
	} else if update.CallbackQuery.Id != "" {
		bot.CallbackQueryHandler()
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (b *Bot) CallbackQueryHandler() {
	callback := b.update.CallbackQuery
	conn := model.Conn{
		DB: b.handler.DB,
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

	token := user.Spotify.Token.User
	newToken, err := spotify.RefreshToken(user.Spotify.Token.Refresh, user.Spotify.Token.Expired)
	if err == nil {
		log.Infof("refresh user token u_id=%s", user.Id)
		if err := conn.UpdateSpotifyToken(&user.Id, &model.Spotify{Token: &model.SpotifyToken{
			Refresh: user.Spotify.Token.Refresh,
			User:    newToken.AccessToken,
			Expired: newToken.Expired,
		}}); err != nil {
			log.Errorf("error update token u_id=%s err=%v", user.Id, err)
			return
		}
		token = newToken.AccessToken
	} else if !errors.Is(err, spotify.TokenNotExpired) {
		log.Errorf("error refresh token u_id=%s err=%v", user.Id, err)
		return
	}

	api := spotify.ApiClient.SetUserToken(token)

	splitData := strings.SplitN(callback.Data, "::", 2)

	if len(splitData) != 2 {
		log.Info("skip ist=", splitData, ", len=", len(splitData))
		return
	}

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
			log.Infof("play song error=%s", err)
		}

		//_ = api.AddQueue(splitData[1])
	} else if splitData[0] == "ADD" {
		if err := api.AddQueue(splitData[1]); err != nil {
			log.Infof("add queue error=%s", err)
		}
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
func (b *Bot) InlineQueryHandler() {
	conn := model.Conn{
		DB: b.handler.DB,
	}

	user, err := GetOrCreateUser(&b.update.InlineQuery.From, b.handler)

	if err != nil {
		return
	}
	if user.Spotify.Token != nil {
		token := user.Spotify.Token.User
		newToken, err := spotify.RefreshToken(user.Spotify.Token.Refresh, user.Spotify.Token.Expired)
		if err == nil {
			log.Infof("refresh user token u_id=%s", user.Id)
			if err := conn.UpdateSpotifyToken(&user.Id, &model.Spotify{Token: &model.SpotifyToken{
				Refresh: user.Spotify.Token.Refresh,
				User:    newToken.AccessToken,
				Expired: newToken.Expired,
			}}); err != nil {
				log.Errorf("error update token u_id=%s err=%v", user.Id, err)
				return
			}
			token = newToken.AccessToken
		} else if !errors.Is(err, spotify.TokenNotExpired) {
			log.Errorf("error refresh token u_id=%s err=%v", user.Id, err)
			return
		}

		api := spotify.ApiClient.SetUserToken(token)
		historyList, err := api.GetAllHistory()
		if err != nil {
			// нужно бы отдавать ошибку в telegram callback
			log.Errorf("error GetAllHistory token u_id=%s err=%v", user.Id, err)
			return
		}

		//tmpList := makePhotoInline(historyList)

		tmpList := makeAudioInline(historyList)

		err = telegram.Client.AnswerInlineQuery(b.update.InlineQuery.Id, tmpList)
		if err != nil {
			log.Error("AnswerInlineQuery err=", err)
		}
	} else {
		r := map[string]interface{}{
			"is_personal":         true,
			"switch_pm_text":      "login in spotify ...",
			"switch_pm_parameter": "inline_redirect",
		}
		err = telegram.Client.AnswerInlineQueryTmp(b.update.InlineQuery.Id, r)
		log.Info("app")
	}
}

func (b *Bot) CommandHandler() {

	user, err := GetOrCreateUser(&b.update.Message.From, b.handler)
	if err != nil {
		log.Error("error create user err=", err)
		return
	}

	command, err := getCommand(b.update.Message.Text)
	// обрабатываем только команды, если нет то скипаем
	if err != nil {
		log.Info("skip text err=", err)
		return
	}

	if user.Active == false {
		log.Info("skip block user =")
		return
	}
	log.Infof("send message: %s", b.update.Message.Text)

	commands := CommandContext{
		update:  b.update,
		user:    user,
		command: command,
		DB:      b.handler.DB,
	}

	switch command.Name {
	case "start":
		commands.StartCommand()
	case "help":
		commands.Help()
	case "setting":

	case "logout":
		// https://support.spotify.com/us/article/how-to-log-out/
	default:

	}
}
