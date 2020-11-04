package commands

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/skar404/spotify_share/spotify"
)

func CreateToken(clientId, clientSecret string) (string, string, error) {
	urlOAuth := spotify.OAuthClient.GetOAthUrl("session_id")

	//c, _ := spotify.Init(clientId, clientSecret, "http://localhost/spotify", []string{
	//	"user-read-recently-played",
	//	"user-read-currently-playing",
	//	"app-remote-control",
	//	"streaming"})

	//urlOAuth := c.GetAuthorizationUrl("session_id")

	fmt.Println(urlOAuth)
	fmt.Print("Enter text: ")

	reader := bufio.NewReader(os.Stdin)
	codeUrl, _ := reader.ReadString('\n')

	codeUrlParse, _ := url.Parse(strings.ReplaceAll(codeUrl, "\n", ""))

	q := codeUrlParse.Query()
	code := q.Get("code")

	r, err := spotify.OAuthClient.GetAccessOrRefreshToken(code)
	if err != nil {
		return "", "", err
	}

	rf, err := spotify.OAuthClient.RefreshToken(r.RefreshToken)
	if err != nil {
		return "", "", err
	}

	api := spotify.ApiClient.SetUserToken(rf.AccessToken)

	a, _ := api.GetPlayNow()
	_ = a

	return rf.AccessToken, r.RefreshToken, nil
}
