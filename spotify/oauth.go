package spotify

import (
	"net/url"
	"strings"

	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/rhttp"
)

type OAuth struct {
	rhttp.ApiClient

	clientId     string
	clientSecret string

	AuthorizationUrl *url.URL
	AccessTokenUrl   string
	ApiUrl           string

	RedirectUri string
	OAuthScope  string

	userToken string
}

var OAuthClient = InitOAuth

func InitOAuth() OAuth {
	urlAuthorization, _ := url.Parse("https://accounts.spotify.com/authorize")

	return OAuth{
		ApiClient: rhttp.ApiClient{
			Url: "https://accounts.spotify.com/api/token",
			Header: map[string][]string{
				"Content-Type": {"application/x-www-form-urlencoded"},
				"Accept": {
					"application/json", "text/json", "text/javascript", "application/xml", "text/xml",
					"application/x-plist", "application/x-www-form-urlencoded", "text/plain", "text/html",
					"application/xhtml+xml"}},
		},
		AuthorizationUrl: urlAuthorization,
		AccessTokenUrl:   "https://accounts.spotify.com/api/token",
		OAuthScope:       strings.Join(global.AppSpotifyScope, " "),
		clientSecret:     global.ClientSecret,
		clientId:         global.ClientId,
	}
}
