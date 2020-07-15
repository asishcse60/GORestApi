package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBClientStorage struct {
	Client       *mongo.Client
	DatabaseName string
}
