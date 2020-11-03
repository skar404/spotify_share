package libs

import (
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"
)

func TestCreateJWT(t *testing.T) {
	// The magic is here

	type args struct {
		userId int64
	}
	tests := []struct {
		name       string
		args       args
		want       string
		monkeyTime time.Time
	}{{
		"",
		args{
			1,
		},
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1NDk4MzI2OTgsImlzcyI6ImxvZ2luX2JvdCJ9.StAfnAjTwQXTCyDpV6mrWUCDDNH5npJP_2DT1dfwkAo",
		time.Date(2019, 02, 10, 20, 34, 58, 651387237, time.UTC),
	}, {
		"",
		args{
			-2401023,
		},
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjotMjQwMTAyMywiZXhwIjoxNTQ5ODMyNjk4LCJpc3MiOiJsb2dpbl9ib3QifQ.6RLr6vHYW_aQo0SFPASzjj1Gm290zFCMyz8_XgafegE",
		time.Date(2019, 02, 10, 20, 34, 58, 651387237, time.UTC),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.Patch(time.Now, func() time.Time {
				return tt.monkeyTime
			})

			if got := CreateJWT(tt.args.userId); got != tt.want {
				t.Errorf("CreateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeJWT(t *testing.T) {
	type args struct {
		jwtString string
	}
	tests := []struct {
		name       string
		args       args
		want       *JWTUser
		wantErr    bool
		monkeyTime time.Time
	}{{
		"",
		args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjotMjQwMTAyMywiZXhwIjoxNjA0NDI5OTM4LCJpc3MiOiJsb2dpbl9ib3QifQ.aHHdw3XR02FZNbJZh67EXA5M1LXQsNkgQA8d9seAdTU"},
		&JWTUser{
			-2401023,
			jwt.StandardClaims{
				Audience:  "",
				ExpiresAt: 1604429938,
				Id:        "",
				IssuedAt:  0,
				Issuer:    "login_bot",
				NotBefore: 0,
				Subject:   "",
			},
		},
		false,
		time.Date(2020, 11, 03, 18, 58, 58, 651387237, time.UTC),
	}, {
		"",
		args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjotMjQwMTAyMywiZXhwIjoxNjA0NDI5OTM4LCJpc3MiOiJsb2dpbl9ib3QifQ.aHHdw3XR02FZNbJZh67EXA5M1LXQsNkgQA8d9seAdTU"},
		nil,
		true,
		time.Date(2020, 11, 03, 18, 59, 58, 651387237, time.UTC),
	}, {
		"", args{"NOT VALID"}, nil, true, time.Now(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.Patch(time.Now, func() time.Time {
				return tt.monkeyTime
			})

			got, err := DecodeJWT(tt.args.jwtString)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
