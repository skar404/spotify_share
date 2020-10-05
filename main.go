package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/skar404/spotify_share/spotify"
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

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

	fmt.Println("GetAccessOrRefreshToken:", r.AccessToken, r.RefreshToken, err)

	rf, err := c.RefreshToken(r)

	fmt.Println("RefreshToken:", rf.AccessToken, err)
}
