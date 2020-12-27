package model

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Conn struct {
	DB        *mongo.Database
	DBMongoDB *mongo.Client
}
