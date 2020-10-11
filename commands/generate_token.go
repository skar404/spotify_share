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
	c, _ := spotify.Init(clientId, clientSecret, "http://localhost/spotify", []string{
		"user-read-recently-played",
		"user-read-currently-playing",
		"app-remote-control",
		"streaming"})

	urlOAuth := c.GetAuthorizationUrl("session_id")

	fmt.Println(urlOAuth)
	fmt.Print("Enter text: ")

	reader := bufio.NewReader(os.Stdin)
	codeUrl, _ := reader.ReadString('\n')

	codeUrlParse, _ := url.Parse(strings.ReplaceAll(codeUrl, "\n", ""))

	q := codeUrlParse.Query()
	code := q.Get("code")

	r, err := c.GetAccessOrRefreshToken(code)
	if err != nil {
		return "", "", err
	}

	rf, err := c.RefreshToken(r)
	if err != nil {
		return "", "", err
	}

	return rf.AccessToken, r.RefreshToken, nil
}
