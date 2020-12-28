package spotify

import (
	"errors"
	"fmt"

	"github.com/skar404/spotify_share/rhttp"
	"github.com/skar404/spotify_share/spotify/type"
)

type api struct {
	rhttp.ApiClient

	userToken string
}

var ApiClient = InitApi()

func InitApi() api {
	return api{
		ApiClient: rhttp.ApiClient{
			Url: "https://api.spotify.com/",
		},
	}
}

var NotFoundError = errors.New("not found")
var NotValidTokenError = errors.New("not valid token")

type QueryError struct {
	Query string
	Err   error
}

func (c api) SetUserToken(token string) api {
	c.userToken = token

	c.ApiClient.Header = map[string][]string{
		"authorization": {fmt.Sprintf("Bearer %s", token)},
	}
	return c
}

// GetPlayNow
// docs api:
// https://developer.spotify.com/documentation/web-api/reference/player/get-the-users-currently-playing-track/
func (c *api) GetPlayNow() (spotify_type.CurrentlyPlaying, error) {
	r := spotify_type.CurrentlyPlaying{}
	response, err := c.HttpClient("GET", "v1/me/player/currently-playing", nil, nil, &r, nil)

	if err != nil {
		return r, NotValidTokenError
	}

	if response.Code != 200 {
		return r, NotFoundError
	}
	return r, nil
}

// GetHistory
// docs api:
// https://developer.spotify.com/documentation/web-api/reference/player/get-recently-played/
func (c *api) GetHistory() (spotify_type.RecentlyPlayed, error) {
	r := spotify_type.RecentlyPlayed{}
	response, err := c.HttpClient("GET", "v1/me/player/recently-played?limit=47", nil, nil, &r, nil)
	if err != nil {
		return r, NotValidTokenError
	}

	if response.Code != 200 {
		return r, NotFoundError
	}
	return r, nil
}

func (c *api) Play(spotifyUri string) error {
	rawData := map[string]interface{}{}
	rawData["uris"] = []string{
		spotifyUri,
	}

	//if contextUri != "" {
	//	//rawData["context_uri"] = contextUri
	//
	//	rawData["offset"] = map[string]string{
	//		"uri": spotifyUri,
	//	}
	//}

	r, err := c.HttpClient("PUT", "v1/me/player/play", rawData, nil, nil, nil)
	_, _ = r, err
	return nil
}

func (c *api) AddQueue(spotifyUri string) error {
	r, err := c.HttpClient("POST", fmt.Sprintf("v1/me/player/queue?uri=%s", spotifyUri), nil, nil, nil, nil)
	_, _ = r, err
	return nil
}

func (c *api) GetPlayer() (spotify_type.Player, error) {
	r := spotify_type.Player{}
	raw, err := c.HttpClient("GET", "v1/me/player", nil, nil, &r, nil)
	_ = raw
	return r, err
}

func (c *api) Next() error {
	r, err := c.HttpClient("POST", "v1/me/player/next", nil, nil, nil, nil)
	_, _ = r, err
	return err
}
