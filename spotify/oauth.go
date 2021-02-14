package spotify

import (
	"github.com/skar404/spotify_share/requests"
	"net/http"
	"net/url"
	"strings"

	"github.com/skar404/spotify_share/global"
	"github.com/skar404/spotify_share/spotify/type"
)

type OAuth struct {
	requests.RequestClient

	clientId     string
	clientSecret string

	AuthorizationUrl *url.URL
	AccessTokenUrl   string
	ApiUrl           string

	RedirectUri string
	OAuthScope  string

	userToken string
}

var OAuthClient = InitOAuth()

func InitOAuth() OAuth {
	urlAuthorization, _ := url.Parse("https://accounts.spotify.com/authorize")

	return OAuth{
		RequestClient: requests.RequestClient{
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
		RedirectUri:      global.RedirectUri,
		clientId:         global.ClientId,
	}
}

func (c *OAuth) GetOAthUrl(state string) string {
	q := c.AuthorizationUrl.Query()

	q.Set("client_id", c.clientId)
	q.Set("response_type", "code")
	q.Set("redirect_uri", c.RedirectUri)
	q.Set("scope", c.OAuthScope)
	q.Set("state", state)
	c.AuthorizationUrl.RawQuery = q.Encode()
	return c.AuthorizationUrl.String()
}

func (c *OAuth) GetAccessOrRefreshToken(code string) (spotify_type.TokenOrRefreshReq, error) {
	var r spotify_type.TokenOrRefreshReq

	req := requests.Request{
		Method: http.MethodPost,
		Uri:    "",
		UrlValues: url.Values{
			"code":          {code},
			"scope":         {c.OAuthScope},
			"redirect_uri":  {c.RedirectUri},
			"client_id":     {c.clientId},
			"client_secret": {c.clientSecret},
			"grant_type":    {"authorization_code"},
		},
	}
	res := requests.Response{
		Struct: &r,
	}
	err := c.NewRequest(&req, &res)
	return r, err
}

func (c *OAuth) RefreshToken(token string) (spotify_type.TokenReq, error) {
	var r spotify_type.TokenReq
	req := requests.Request{
		Method: http.MethodPost,
		Uri:    "",
		UrlValues: url.Values{
			"scope":         {c.OAuthScope},
			"refresh_token": {token},
			"client_id":     {c.clientId},
			"client_secret": {c.clientSecret},
			"grant_type":    {"refresh_token"},
		},
	}
	res := requests.Response{
		Struct: &r,
	}
	err := c.NewRequest(&req, &res)
	return r, err
}
