package spotify

import (
	"net/url"
	"strings"

	"github.com/skar404/spotify_share/rhttp"
)

type Config struct {
	// FIXME возможно нужно разбить на 2 клиента
	rhttp.ApiClient
	AuthorizationClient rhttp.ApiClient

	clientId     string
	clientSecret string

	AuthorizationUrl *url.URL
	AccessTokenUrl   string
	ApiUrl           string

	RedirectUri string
	OAuthScope  string

	userToken string
}

func Init(id string, secret string, redirectUri string, scope []string) (Config, error) {
	urlAuthorization, err := url.Parse("https://accounts.spotify.com/authorize")

	if err != nil {
		return Config{}, err
	}

	return Config{
		ApiClient: rhttp.ApiClient{
			Url: "https://api.spotify.com/",
		},
		AuthorizationClient: rhttp.ApiClient{
			Url: "https://accounts.spotify.com/api/token",
			Header: map[string][]string{
				"Content-Type": {"application/x-www-form-urlencoded"},
				"Accept": {
					"application/json", "text/json", "text/javascript", "application/xml", "text/xml",
					"application/x-plist", "application/x-www-form-urlencoded", "text/plain", "text/html",
					"application/xhtml+xml"}},
		},

		clientId:     id,
		clientSecret: secret,

		AuthorizationUrl: urlAuthorization,

		// FIXME не уверне что нужно дублироваь AccessTokenUrl и ApiUrl
		AccessTokenUrl: "https://accounts.spotify.com/api/token",
		ApiUrl:         "https://api.spotify.com/",

		OAuthScope: strings.Join(scope, " "),

		RedirectUri: redirectUri,
	}, nil
}

func (c *Config) GetAuthorizationUrl(state string) string {
	q := c.AuthorizationUrl.Query()

	q.Set("client_id", c.clientId)
	q.Set("response_type", "code")
	q.Set("redirect_uri", c.RedirectUri)
	q.Set("scope", c.OAuthScope)
	q.Set("state", state)
	c.AuthorizationUrl.RawQuery = q.Encode()
	return c.AuthorizationUrl.String()
}

func (c *Config) GetAccessOrRefreshToken(code string) (TokenOrRefreshReq, error) {
	var r TokenOrRefreshReq
	var err error

	body := url.Values{
		"code":          {code},
		"scope":         {c.OAuthScope},
		"redirect_uri":  {c.RedirectUri},
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
		"grant_type":    {"authorization_code"},
	}

	_, err = c.AuthorizationClient.HttpClient("POST", "", nil, body, &r, nil)
	return r, err
}

func (c *Config) RefreshToken(token TokenOrRefreshReq) (TokenReq, error) {
	var r TokenReq
	var err error

	body := url.Values{
		"scope":         {c.OAuthScope},
		"refresh_token": {token.RefreshToken},
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
		"grant_type":    {"refresh_token"},
	}

	_, err = c.AuthorizationClient.HttpClient("POST", "", nil, body, &r, nil)

	return r, err
}
