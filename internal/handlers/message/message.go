package messagehandlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	messagedata "messaging/internal/data/message"
	shareddata "messaging/internal/data/shared"
	sharedhandlers "messaging/internal/handlers/shared"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Messages struct {
	l      *log.Logger
	client *mongo.Client
	v      *shareddata.Validation
}

func NewMessages(l *log.Logger, c *mongo.Client, v *shareddata.Validation) *Messages {
	return &Messages{l: l, client: c, v: v}
}

type KeyMessages struct{}

func (m *Messages) MiddlewareMessagesValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		msg := &messagedata.Message{}
		err := shareddata.FromJSON(msg, r.Body)

		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			m.l.Printf("Error unmarshaling messages: %s\n", err)
			return
		}

		errs := m.v.Validate(msg)

		if len(errs) != 0 {
			m.l.Printf("[ERROR] validating message: %+v\n", errs)
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			shareddata.ToJSON(&sharedhandlers.ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyMessages{}, msg)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)

	})
}

func (m *Messages) GetMessage(rw http.ResponseWriter, r *http.Request) {
	type req struct {
		Bus              string `json:"bus"`
		NumberOfMessages int    `json:"numberOfMessages"`
	}

	reqBody := req{}

	read, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(read, &reqBody)
	if err != nil {
		m.l.Printf("Error unmarshaling request body: %s\n", err)
		http.Error(rw, "Unable to marshal request body", http.StatusBadRequest)
		return
	}

	res, err := messagedata.ReceiveMessages(r.Context(), messagedata.ReceiveMessagesInput{
		Client:           m.client,
		NumberOfMessages: reqBody.NumberOfMessages,
		Bus:              reqBody.Bus,
	})

	resBody, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	rw.Write(resBody)
}

func (m *Messages) PostMessage(rw http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value(KeyMessages{}).(*messagedata.Message)

	_, err := messagedata.SendMessage(r.Context(), messagedata.SendMessagesInput{
		Client:  m.client,
		Message: msg,
	})

	if err != nil {
		http.Error(rw, "Unable to send messages to database", http.StatusInternalServerError)
		m.l.Printf("Error sending messages: %s\n", err)
	}
}
