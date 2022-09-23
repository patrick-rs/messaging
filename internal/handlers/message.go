package handlers

import (
	"log"
	"messaging/internal/data"

	"go.mongodb.org/mongo-driver/mongo"
)

type Messages struct {
	l      *log.Logger
	client *mongo.Client
	v      *data.Validation
}

type Bus struct {
	l      *log.Logger
	client *mongo.Client
	v      *data.Validation
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

func NewMessages(l *log.Logger, c *mongo.Client, v *data.Validation) *Messages {
	return &Messages{l: l, client: c, v: v}
}

func NewBus(l *log.Logger, c *mongo.Client, v *data.Validation) *Bus {
	return &Bus{l: l, client: c, v: v}
}
