package busdata

import (
	"context"
	"encoding/json"
	"io"
	datashared "messaging/internal/data/shared"

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

type GetBusInput struct {
	Client *mongo.Client
	Bus    *Bus
}

type GetBusOutput struct {
	Bus *Bus
}

func GetBus(ctx context.Context, in GetBusInput) (GetBusOutput, error) {
	out := GetBusOutput{}

	coll := in.Client.Database(datashared.DATABASE).Collection(datashared.BUSES_COLLECTION)

	bus := Bus{}

	query := bson.M{"name": in.Bus.Name}

	err := coll.FindOne(ctx, query).Decode(&bus)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return out, nil
		}
		return out, err
	}

	out.Bus = &bus

	return out, nil
}
