package bushandlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	busdata "messaging/internal/data/bus"
	shareddata "messaging/internal/data/shared"
	sharedhandlers "messaging/internal/handlers/shared"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewBus(l *log.Logger, c *mongo.Client, v *shareddata.Validation) *Bus {
	return &Bus{l: l, client: c, v: v}
}

type Bus struct {
	l      *log.Logger
	client *mongo.Client
	v      *shareddata.Validation
}
type KeyBus struct{}

func (b *Bus) MiddlewareBusesValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bus := &busdata.Bus{}
		err := shareddata.FromJSON(bus, r.Body)

		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			b.l.Printf("Error unmarshaling messages: %s\n", err)
			return
		}

		errs := b.v.Validate(bus)

		if len(errs) != 0 {
			b.l.Printf("[ERROR] validating message: %+v\n", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			shareddata.ToJSON(&sharedhandlers.ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyBus{}, bus)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)

	})
}

func (b *Bus) GetBus(rw http.ResponseWriter, r *http.Request) {

	bus := busdata.Bus{}

	read, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(read, &bus)

	if err != nil {
		b.l.Printf("Error unmarshaling request body: %s\n", err)
		http.Error(rw, "Unable to marshal request body", http.StatusBadRequest)
		return
	}

	res, err := busdata.GetBus(r.Context(), busdata.GetBusInput{
		Client: b.client,
		Bus:    &bus,
	})

	if err != nil {
		b.l.Printf("Error unmarshaling bus: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(res)

	if err != nil {
		b.l.Printf("Error unmarshaling bus: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(resBody)
}

func (b *Bus) PostBus(rw http.ResponseWriter, r *http.Request) {

	bus := busdata.Bus{}

	read, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(read, &bus)

	if err != nil {
		b.l.Printf("Error unmarshaling request body: %s\n", err)
		http.Error(rw, "Unable to marshal request body", http.StatusBadRequest)
		return
	}

	err = busdata.CreateBus(r.Context(), busdata.CreateBusInput{
		Client: b.client,
		Bus:    &bus,
	})

	if err != nil {
		b.l.Printf("Error creating bus: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
}

func (b *Bus) DeleteBus(rw http.ResponseWriter, r *http.Request) {

	bus := busdata.Bus{}

	read, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(read, &bus)

	if err != nil {
		b.l.Printf("Error unmarshaling request body: %s\n", err)
		http.Error(rw, "Unable to marshal request body", http.StatusBadRequest)
		return
	}

	err = busdata.CreateBus(r.Context(), busdata.CreateBusInput{
		Client: b.client,
		Bus:    &bus,
	})

	if err != nil {
		b.l.Printf("Error creating bus: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
}
