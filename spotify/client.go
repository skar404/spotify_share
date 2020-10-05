package spotify

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/skar404/spotify_share/http"
)

type Config struct {
	http.ApiClient
	AuthorizationClient http.ApiClient

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
		ApiClient: http.ApiClient{
			Url: "https://api.spotify.com/",
		},
		AuthorizationClient: http.ApiClient{
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

func (c *Config) GetAccessOrRefreshToken(code string) {
	body := url.Values{
		"code":          {code},
		"scope":         {code},
		"redirect_uri":  {c.RedirectUri},
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
		"grant_type":    {"authorization_code"},
	}

	r, _ := c.AuthorizationClient.HttpClient("POST", "", nil, body, nil, nil)

	fmt.Println(r)
}
