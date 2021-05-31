package spotify

import (
	"errors"
	"fmt"
	"github.com/skar404/spotify_share/requests"
	"github.com/skar404/spotify_share/spotify/type"
	"net/http"
)

type ApiContext struct {
	requests.RequestClient

	userToken string
}

var ApiClient = InitApi()

func InitApi() ApiContext {
	return ApiContext{
		RequestClient: requests.RequestClient{
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

func (c ApiContext) SetUserToken(token string) ApiContext {
	c.userToken = token
	c.RequestClient.Header = map[string][]string{
		"authorization": {fmt.Sprintf("Bearer %s", token)},
	}
	return c
}

// GetPlayNow
// docs ApiContext:
// https://developer.spotify.com/documentation/web-api/reference/player/get-the-users-currently-playing-track/
func (c *ApiContext) GetPlayNow() (spotify_type.CurrentlyPlaying, error) {
	r := spotify_type.CurrentlyPlaying{}

	req := requests.Request{
		Method: http.MethodGet,
		Uri:    "v1/me/player/currently-playing",
	}
	res := requests.Response{
		Struct: &r,
	}
	err := c.NewRequest(&req, &res)
	if res.Code == http.StatusNoContent {
		return r, NotFoundError
	}
	if err != nil {
		return r, NotValidTokenError
	}

	// If the device is not found, the request will return 404 NOT FOUND response code.
	// If the user making the request is non-premium, a 403 FORBIDDEN response code will be returned.
	if res.Code != http.StatusOK {
		return r, NotFoundError
	}
	return r, nil
}

// GetHistory
// docs ApiContext:
// https://developer.spotify.com/documentation/web-api/reference/player/get-recently-played/
func (c *ApiContext) GetHistory(limit int) (spotify_type.RecentlyPlayed, error) {
	r := spotify_type.RecentlyPlayed{}
	req := requests.Request{
		Method: http.MethodGet,
		Uri:    fmt.Sprintf("v1/me/player/recently-played?limit=%d", limit),
	}
	res := requests.Response{
		Struct: &r,
	}
	err := c.NewRequest(&req, &res)
	if err != nil {
		return r, NotValidTokenError
	}
	if res.Code != http.StatusOK {
		return r, NotFoundError
	}
	return r, nil
}

func (c *ApiContext) Play(spotifyUri string) error {
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

	req := requests.Request{
		Method:   http.MethodPut,
		Uri:      "v1/me/player/play",
		JsonBody: rawData,
	}
	res := requests.Response{}

	err := c.NewRequest(&req, &res)
	if err != nil {
		return NotValidTokenError
	}
	if res.Code != http.StatusNoContent {
		return NotFoundError
	}

	return nil
}

func (c *ApiContext) AddQueue(spotifyUri string) error {
	req := requests.Request{
		Method: http.MethodPost,
		Uri:    fmt.Sprintf("v1/me/player/queue?uri=%s", spotifyUri),
	}
	res := requests.Response{}

	err := c.NewRequest(&req, &res)
	if err != nil {
		return NotValidTokenError
	}
	if res.Code != http.StatusNoContent {
		return NotFoundError
	}
	return nil
}

func (c *ApiContext) AddTracks(ID string) error {
	req := requests.Request{
		Method: http.MethodPut,
		Uri:    "v1/me/tracks",
		JsonBody: map[string][]string{
			"ids": {ID},
		},
	}
	res := requests.Response{}

	err := c.NewRequest(&req, &res)
	if err != nil {
		return NotValidTokenError
	}
	if res.Code != http.StatusOK {
		return fmt.Errorf("res=%+v err=%s", res, NotFoundError)
	}
	return nil
}

func (c *ApiContext) GetPlayer() (spotify_type.Player, error) {
	r := spotify_type.Player{}

	req := requests.Request{
		Method: http.MethodGet,
		Uri:    "v1/me/player",
	}
	res := requests.Response{
		Struct: &r,
	}

	err := c.NewRequest(&req, &res)
	if err != nil {
		return r, NotValidTokenError
	}
	if res.Code != http.StatusOK {
		return r, NotFoundError
	}

	return r, err
}

func (c *ApiContext) Next() error {
	req := requests.Request{
		Method: http.MethodPost,
		Uri:    "v1/me/player/next",
	}
	res := requests.Response{}

	err := c.NewRequest(&req, &res)
	if err != nil {
		return NotValidTokenError
	}
	if res.Code != http.StatusOK {
		return NotFoundError
	}
	return nil
}
