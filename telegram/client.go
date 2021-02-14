package telegram

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/requests"
)

type Context struct {
	requests.RequestClient
}

var Client = Init()

func Init() Context {
	return Context{
		requests.RequestClient{
			Url: fmt.Sprintf("https://api.telegram.org/bot%s/", global.TelegramToken),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Timeout: 1 * time.Minute,
		},
	}
}

func (c *Context) SendMessage(chatId int64, text string, replyMarkup *InlineKeyboardReq) error {
	jsonBody := map[string]interface{}{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "Markdown",
	}

	if replyMarkup != nil && replyMarkup.ready == true {
		jsonBody["reply_markup"] = replyMarkup
	}

	req := requests.Request{
		Method:   http.MethodPost,
		Uri:      "sendMessage",
		JsonBody: &jsonBody,
	}
	res := requests.Response{}
	return c.NewRequest(&req, &res)
}

func (c *Context) GetUpdates(offSet int) (*GetUpdate, error) {
	jsonBody := make(map[string]interface{})
	if offSet != 0 {
		jsonBody["offset"] = strconv.Itoa(offSet)
		jsonBody["timeout"] = 1
	}

	r := &GetUpdate{}
	req := requests.Request{
		Method:   http.MethodPost,
		Uri:      "getUpdates",
		JsonBody: &jsonBody,
	}
	res := requests.Response{
		Struct: &r,
	}
	err := Client.NewRequest(&req, &res)
	return r, err
}

func (c *Context) AnswerInlineQuery(Id string, tmpList []interface{}) error {
	jsonBody := make(map[string]interface{})

	jsonBody["inline_query_id"] = Id
	jsonBody["cache_time"] = 0
	jsonBody["results"] = tmpList

	r := &GetUpdate{}
	req := requests.Request{
		Method:   http.MethodPost,
		Uri:      "answerInlineQuery",
		JsonBody: &jsonBody,
	}
	res := requests.Response{
		Struct: &r,
	}

	return c.NewRequest(&req, &res)
}

func (c *Context) AnswerInlineQueryTmp(Id string, jsonBody map[string]interface{}) error {
	jsonBody["inline_query_id"] = Id

	r := &GetUpdate{}
	req := requests.Request{
		Method:   http.MethodPost,
		Uri:      "answerInlineQuery",
		JsonBody: &jsonBody,
	}
	res := requests.Response{
		Struct: &r,
	}
	return c.NewRequest(&req, &res)
}

func (c *Context) AnswerCallbackQuery(Id string, data *AnswerCallbackReq) error {
	rawData := map[string]interface{}{}
	rawData["callback_query_id"] = Id
	rawData["text"] = data.Text
	rawData["show_alert"] = data.ShowAlert
	rawData["url"] = data.Url

	r := &GetUpdate{}
	req := requests.Request{
		Method:   http.MethodPost,
		Uri:      "answerCallbackQuery",
		JsonBody: &rawData,
	}
	res := requests.Response{
		Struct: &r,
	}
	return c.NewRequest(&req, &res)
}
