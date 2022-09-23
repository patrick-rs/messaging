package handlers

import (
	"context"
	"messaging/internal/data"
	"net/http"
)

type KeyMessages struct{}

func (m *Messages) MiddlewareMessagesValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		msg := &data.Message{}
		err := data.FromJSON(msg, r.Body)

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
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyMessages{}, msg)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)

	})
}

type KeyBus struct{}

func (b *Bus) MiddlewareBusesValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bus := &data.Bus{}
		err := data.FromJSON(bus, r.Body)

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
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyBus{}, bus)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)

	})
}
