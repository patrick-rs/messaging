package shareddata

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DATABASE            = "messaging"
	MESSAGES_COLLECTION = "messages"
	BUSES_COLLECTION    = "buses"
)

func NewMongoDBClient() (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://db:%s/", os.Getenv("DB_PORT"))
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
