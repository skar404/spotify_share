package bot

import (
	"errors"
	"fmt"
	"github.com/skar404/spotify_share/spotify"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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

func GetOrCreateUser(tgUser *telegram.User, h *handler.Handler) (*model.User, error) {
	if tgUser.Id == 0 {
		return nil, NotValidUser
	}

	conn := model.Conn{
		DB: h.DB,
	}
	user, err := conn.GetUser(tgUser.Id)

	if err == nil {
		return user, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	u := &model.User{
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

func makePhotoInline(h []spotify.History) []interface{} {
	// FIXME перепесать это де**** на struct
	//  p.s. автор данного куса не несет ответсвенность за ваше психическое состояние

	tmpList := make([]interface{}, len(h))
	for i := range h {
		link := &h[i]

		title := link.Name
		if link.PlayNow {
			title = "▶️ ️" + title
		}

		tmpList[i] = map[string]interface{}{
			"type":  "photo",
			"id":    fmt.Sprintf("%v %v", time.Now().Unix(), RandStringBytes(10)),
			"title": title,
			"description": fmt.Sprintf("%s",
				link.Artists[0].Name),
			"caption": fmt.Sprintf("Name: ***%s***\nArtist: ***%s***\nAlbum: ***%s***",
				link.Name, link.Artists[0].Name, link.Album.Name),
			"photo_url":  link.Img,
			"parse_mode": "Markdown",
			"reply_markup": map[string][][]map[string]string{
				"inline_keyboard": {{
					{
						"text":          "Play",
						"callback_data": fmt.Sprintf("PLAY::%s", link.URL),
					},
					{
						"text":          "Add",
						"callback_data": fmt.Sprintf("ADD::%s", link.URL),
					},
				}},
			},
			"thumb_url": link.Img,
		}
	}
	return tmpList
}

func makeAudioInline(h []spotify.History) []interface{} {
	// FIXME перепесать это де**** на struct
	//  p.s. автор данного куса не несет ответсвенность за ваше психическое состояние

	tmpList := make([]interface{}, len(h))
	for i := range h {
		link := &h[i]

		title := link.Name
		if link.PlayNow {
			title = "▶️ ️" + title
		}

		tmpList[i] = map[string]interface{}{
			"type":           "audio",
			"id":             fmt.Sprintf("%v %v", time.Now().Unix(), RandStringBytes(10)),
			"audio_url":      link.PreviewURL,
			"title":          title,
			"caption":        fmt.Sprintf("[song link](https://song.link/s/%s)", link.URL),
			"parse_mode":     "Markdown",
			"performer":      link.Artists[0].Name,
			"audio_duration": 30, // вроде по всем трекам отдает трек 30c, если найду меньше нужно потестит, НО это не критично
			// идея скачивать трек ОЧЕНЬ ПЛОХАЯ
			"reply_markup": map[string][][]map[string]string{
				"inline_keyboard": {{
					{
						"text":          "Play",
						"callback_data": fmt.Sprintf("PLAY::%s", link.URL),
					},
					{
						"text":          "Add",
						"callback_data": fmt.Sprintf("ADD::%s", link.URL),
					},
					//{
					//	"text":          "Like",
					//	"callback_data": fmt.Sprintf("LIKE::%s", link.URL),
					//},
				}},
			},
			"thumb_url": link.Img,
		}
	}
	return tmpList
}
