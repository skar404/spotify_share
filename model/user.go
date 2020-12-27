package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

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
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	Telegram Telegram `json:"telegram,omitempty" bson:"telegram,omitempty"`

	Spotify Spotify `json:"spotify,omitempty" bson:"spotify,omitempty"`

	CreateAt time.Time `json:"create_at,omitempty" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty" bson:"update_at,omitempty"`

	Active bool `json:"active,omitempty" bson:"active,omitempty"`
}

func (c *Conn) Collection() *mongo.Collection {
	coll := c.DB.Collection(userCollection)
	return coll
}

func (c *Conn) GetUser(tgId int64) (*User, error) {
	var err error
	ctx := context.Background()
	collection := c.Collection()

	res := collection.FindOne(ctx, bson.M{"telegram.id": tgId})

	findUser := User{}
	err = res.Decode(&findUser)

	if err != nil {
		return nil, err
	}

	return &findUser, nil
}

// CreateUser
// ... данный мметод принемает и отдает одну и тожу ссылку,
// ... но сделал это для удобства и большей явности
func (c *Conn) CreateUser(u *User) (*User, error) {
	var err error
	ctx := context.Background()
	collection := c.Collection()

	result, err := collection.InsertOne(ctx, u)
	if err != nil {
		return u, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		u.Id = oid
	} else {
		// FIXME add context in err
		return u, errors.New("error get ObjectID")
	}

	return u, nil
}

func (c *Conn) UpdateSpotifyToken(u *primitive.ObjectID, sToken *Spotify) error {
	ctx := context.Background()
	collection := c.Collection()
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": &u},
		bson.D{
			{"$set", map[string]time.Time{"update_at": time.Now()}},
			{"$set", map[string]Spotify{
				"spotify": *sToken,
			}},
		},
	)
	return err
}
