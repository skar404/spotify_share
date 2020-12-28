package telegram

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/rhttp"
)

type Config struct {
	rhttp.ApiClient

	Token string
}

var TgClient = Init()

func Init() Config {
	token := global.TelegramToken
	return Config{
		ApiClient: rhttp.ApiClient{
			Url: fmt.Sprintf("https://api.telegram.org/bot%s/", token),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *Config) SendMessage(chatId int64, text string, replyMarkup *InlineKeyboardReq) error {
	jsonBody := map[string]interface{}{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "Markdown",
	}

	if replyMarkup != nil && replyMarkup.ready == true {
		jsonBody["reply_markup"] = replyMarkup
	}

	r, err := c.HttpClient("POST", "sendMessage", jsonBody, nil, nil, nil)
	_ = r
	return err
}

func (c *Config) SetChatDescription(chatId int, text string) error {
	jsonBody := map[string]interface{}{
		"chat_id":     chatId,
		"description": text,
	}

	_, err := c.HttpClient("POST", "setChatDescription", jsonBody, nil, nil, nil)
	return err
}

func (c *Config) GetMe() (GetMe, error) {
	resUpdate := GetMe{}

	_, err := c.HttpClient("GET", "getMe", nil, nil, &resUpdate, nil)
	return resUpdate, err
}

func (c *Config) SetWebHook(hookUrl string, maxConn int) error {
	jsonBody := map[string]interface{}{
		"url": hookUrl,
	}

	if maxConn != 0 {
		jsonBody["max_connections"] = maxConn
	}

	_, err := c.HttpClient("POST", "setWebhook", jsonBody, nil, nil, nil)
	return err
}

func (c *Config) GetUpdates(offSet int) (*GetUpdate, error) {
	jsonBody := make(map[string]interface{})
	if offSet != 0 {
		jsonBody["offset"] = strconv.Itoa(offSet)
		jsonBody["timeout"] = 30
	}

	resUpdate := &GetUpdate{}
	r, err := c.HttpClient("POST", "getUpdates", jsonBody, nil, &resUpdate, nil)
	_ = r
	return resUpdate, err
}

func (c *Config) AnswerInlineQuery(Id string, tmpList []interface{}) error {

	jsonBody := make(map[string]interface{})

	jsonBody["inline_query_id"] = Id
	jsonBody["cache_time"] = 0
	jsonBody["results"] = tmpList

	resUpdate := GetUpdate{}
	r, err := c.HttpClient("POST", "answerInlineQuery", jsonBody, nil, &resUpdate, nil)
	_ = r
	return err
}

func (c *Config) AnswerInlineQueryTmp(Id string, jsonBody map[string]interface{}) error {
	jsonBody["inline_query_id"] = Id

	resUpdate := GetUpdate{}
	r, err := c.HttpClient("POST", "answerInlineQuery", jsonBody, nil, &resUpdate, nil)
	_ = r

	log.Info(r.Raw)
	return err
}

func (c *Config) AnswerCallbackQuery(Id string, data *AnswerCallbackReq) error {
	rawData := map[string]interface{}{}
	rawData["callback_query_id"] = Id
	rawData["text"] = data.Text
	rawData["show_alert"] = data.ShowAlert
	rawData["url"] = data.Url

	r, err := c.HttpClient("POST", "answerCallbackQuery", rawData, nil, nil, nil)
	_ = r
	return err
}
