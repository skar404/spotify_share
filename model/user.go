package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Telegram struct {
		Id    int64  `json:"id,omitempty" bson:"id,omitempty"`
		Login string `json:"login" bson:"login"`
	} `json:"telegram,omitempty" bson:"telegram,omitempty"`

	Spotify struct {
		Token struct {
			Refresh string `json:"refresh" bson:"refresh"`
			User    string `json:"user" bson:"user"`
		} `json:"token,omitempty" bson:"token,omitempty"`
	} `json:"spotify,omitempty" bson:"spotify,omitempty"`

	CreateAt time.Time `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`
}
