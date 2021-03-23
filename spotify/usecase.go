package spotify

import (
	"errors"
	"fmt"
	"time"
)

type RefreshTokenRes struct {
	AccessToken string
	Expired     int64
}

var TokenNotExpired = errors.New("token not expired")

func RefreshToken(refreshToken string, expired int64) (RefreshTokenRes, error) {
	var r RefreshTokenRes

	timeStamp := time.Now().Unix()
	if expired > timeStamp {
		return r, TokenNotExpired
	}

	token, err := OAuthClient.RefreshToken(refreshToken)
	if err != nil {
		return RefreshTokenRes{}, fmt.Errorf("failed refresh token err=%e", err)
	}

	return RefreshTokenRes{
		AccessToken: token.AccessToken,
		Expired:     timeStamp + token.ExpiresIn,
	}, nil
}
