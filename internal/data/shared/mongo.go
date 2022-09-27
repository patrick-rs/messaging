package shareddata

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	uri                 = "mongodb://localhost:27017"
	DATABASE            = "messaging"
	MESSAGES_COLLECTION = "messages"
	BUSES_COLLECTION    = "buses"
)

func NewMongoDBClient() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
