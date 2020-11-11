package model

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const userCollection = "users"

var UserNotActive = errors.New("user not active")

type SpotifyToken struct {
	Refresh string `json:"refresh" bson:"refresh"`
	User    string `json:"user" bson:"user"`
	Expired int64  `json:"expired" bson:"expired"`
}

type Spotify struct {
	Token SpotifyToken `json:"token,omitempty" bson:"token,omitempty"`
}

type Telegram struct {
	Id    int64  `json:"id,omitempty" bson:"id,omitempty"`
	Login string `json:"login" bson:"login"`
}

type User struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Telegram Telegram `json:"telegram,omitempty" bson:"telegram,omitempty"`

	Spotify Spotify `json:"spotify,omitempty" bson:"spotify,omitempty"`

	CreateAt time.Time `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`

	Active bool `json:"active,omitempty" bson:"active,omitempty"`
}

func (c *Conn) Collection() (*mgo.Session, *mgo.Collection) {
	conn, db := c.Clone()
	table := db.C(userCollection)
	return conn, table
}

func (c *Conn) GetUser(tgId int64) (*User, error) {
	var err error
	conn, collection := c.Collection()
	defer conn.Close()

	u := &User{Id: bson.NewObjectId()}

	err = collection.Find(bson.M{"telegram.id": tgId}).One(u)

	if err != nil {
		return nil, err
	}

	if u.Active == false {
		return nil, UserNotActive
	}

	return u, nil
}

// CreateUser
// ... данный мметод принемает и отдает одну и тожу ссылку,
// ... но сделал это для удобства и большей явности
func (c *Conn) CreateUser(u *User) (*User, error) {
	var err error
	conn, collection := c.Collection()
	defer conn.Close()

	err = collection.Insert(u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (c *Conn) UpdateUser(u *bson.ObjectId, user *User) error {
	conn, collection := c.Collection()
	defer conn.Close()

	user.UpdateAt = time.Now()

	err := collection.UpdateId(u, bson.M{"$set": user})
	return err
}
