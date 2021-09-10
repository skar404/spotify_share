package libs

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/skar404/spotify_share/global"
)

type JWTUser struct {
	UserId int64 `json:"user_id,omitempty"`
	jwt.StandardClaims
}

func CreateJWT(userId int64) string {
	mySigningKey := global.JWTTokenByte

	date := time.Now()
	date = date.Add(5 * time.Minute)

	claims := JWTUser{
		userId,
		jwt.StandardClaims{
			ExpiresAt: date.Unix(),
			Issuer:    "login_bot",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return ""
	}
	return ss
}

func DecodeJWT(jwtString string) (*JWTUser, error) {
	token, err := jwt.ParseWithClaims(jwtString, &JWTUser{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JWTToken), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTUser); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
