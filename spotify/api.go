package spotify

import "github.com/skar404/spotify_share/rhttp"

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

func (s *api) SetUserToken(token string) api {
	s.userToken = token

	return *s
}
