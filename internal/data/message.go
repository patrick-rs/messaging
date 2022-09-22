package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Bus  string `bson:"bus" json:"bus" validate:"required"`
	Body string `bson:"body" json:"body" validate:"required"`
}

type Messages []*Message

func (p *Messages) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Messages) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

const (
	database   = "messaging"
	collection = "messages"
)

type SendMessagesInput struct {
	Client  *mongo.Client
	Message *Message
}

type SendMessageOutput struct {
	ID interface{}
}

func SendMessage(ctx context.Context, in SendMessagesInput) (SendMessageOutput, error) {
	coll := in.Client.Database(database).Collection(collection)
	/*
		data := make([]interface{}, len(*in.Messages))
		for i, d := range *in.Messages {
			data[i] = d
		}
	*/
	res, err := coll.InsertOne(ctx, in.Message)
	out := SendMessageOutput{}
	if err != nil {
		return out, err
	}
	out.ID = res.InsertedID
	return out, nil
}

type ReceiveMessagesInput struct {
	Client           *mongo.Client
	NumberOfMessages int
	Bus              string
}

type ReceiveMessagesOutput struct {
	Messages Messages
}

func ReceiveMessages(ctx context.Context, in ReceiveMessagesInput) (ReceiveMessagesOutput, error) {
	coll := in.Client.Database(database).Collection(collection)
	out := ReceiveMessagesOutput{Messages: Messages{}}

	fmt.Println("BUS", in.Bus)
	query := bson.M{"bus": in.Bus}

	cursor, err := coll.Find(ctx, query)

	if err != nil {
		return out, err
	}

	for cursor.Next(ctx) {
		msg := Message{}
		err := cursor.Decode(&msg)

		if err != nil {
			return out, err
		}
		out.Messages = append(out.Messages, &msg)
	}

	return out, nil
}
