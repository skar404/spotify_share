package telegram

import (
	"fmt"
	"strconv"

	"github.com/skar404/spotify_share/rhttp"
)

type Config struct {
	rhttp.ApiClient

	Token string
}

func Init(token string) (Config, error) {
	return Config{
		ApiClient: rhttp.ApiClient{
			Url: fmt.Sprintf("https://api.telegram.org/bot%s/", token),
		},
	}, nil
}

func (c *Config) SendMessage(chatId int, text string) error {
	jsonBody := map[string]interface{}{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "Markdown",
	}

	_, err := c.HttpClient("POST", "sendMessage", jsonBody, nil, nil, nil)
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

func (c *Config) GetUpdates(offSet int) (GetUpdate, error) {

	jsonBody := make(map[string]interface{})
	if offSet != 0 {
		jsonBody["offset"] = strconv.Itoa(offSet)
	}

	resUpdate := GetUpdate{}
	_, err := c.HttpClient("POST", "getUpdates", jsonBody, nil, &resUpdate, nil)

	return resUpdate, err
}
