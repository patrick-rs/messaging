package data

import (
	"context"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bus struct {
	Name string `bson:"name" json:"name" validate:"required"`
}

func (b *Bus) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(b)
}

func (b *Bus) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(b)
}

type CheckIfBusExistsInput struct {
	Client *mongo.Client
	Bus    *Bus
}

type CheckIfBusExistsOutput struct {
	Bus    *Bus
	Exists bool
}

func CheckIfBusExists(ctx context.Context, in CheckIfBusExistsInput) (CheckIfBusExistsOutput, error) {
	out := CheckIfBusExistsOutput{}

	coll := in.Client.Database(database).Collection(BUSES_COLLECTION)

	bus := Bus{}

	query := bson.M{"name": in.Bus.Name}

	err := coll.FindOne(ctx, query).Decode(&bus)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			out.Exists = false
			return out, nil
		}
		return out, err
	}

	out.Exists = true
	out.Bus = &bus

	return out, nil
}
