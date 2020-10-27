package spotify

import (
	"fmt"

	"github.com/skar404/spotify_share/rhttp"
)

type api struct {
	rhttp.ApiClient

	userToken string
}

var ApiClient = InitApi()

func InitApi() api {
	return api{
		ApiClient: rhttp.ApiClient{
			Url:    "https://api.spotify.com/",
			Header: map[string][]string{},
		},
	}
}

func (c *api) SetUserToken(token string) api {
	c.userToken = token

	c.ApiClient.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	return *c
}

func (c *api) GetPlayNow() rhttp.ResponseJson {
	r, _ := c.HttpClient("GET", "v1/me/player/currently-playing", nil, nil, nil, nil)
	return r
}

func (c *api) GetHistory() rhttp.ResponseJson {
	r, _ := c.HttpClient("GET", "v1/me/player/recently-played", nil, nil, nil, nil)
	return r
}
