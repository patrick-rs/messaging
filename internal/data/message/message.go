package datamessage

import (
	"context"
	"encoding/json"
	"io"
	datashared "messaging/internal/data/shared"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Bus  string `bson:"bus" json:"bus" validate:"required"`
	Body string `bson:"body" json:"body" validate:"required"`
}

type Messages []*Message

func (m *Messages) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func (m *Messages) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

type SendMessagesInput struct {
	Client  *mongo.Client
	Message *Message
}

type SendMessageOutput struct {
	ID interface{}
}

func SendMessage(ctx context.Context, in SendMessagesInput) (SendMessageOutput, error) {
	coll := in.Client.Database(datashared.DATABASE).Collection(datashared.MESSAGES_COLLECTION)
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
	coll := in.Client.Database(datashared.DATABASE).Collection(datashared.MESSAGES_COLLECTION)
	out := ReceiveMessagesOutput{Messages: Messages{}}

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
