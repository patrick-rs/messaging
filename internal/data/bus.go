package data

import (
	"encoding/json"
	"io"
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
