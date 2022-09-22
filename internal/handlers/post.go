package handlers

import (
	"messaging/internal/data"
	"net/http"
)

func (m *Messages) PostMessage(rw http.ResponseWriter, r *http.Request) {
	msg := r.Context().Value(KeyMessages{}).(*data.Message)

	_, err := data.SendMessage(r.Context(), data.SendMessagesInput{
		Client:  m.client,
		Message: msg,
	})

	if err != nil {
		http.Error(rw, "Unable to send messages to database", http.StatusInternalServerError)
		m.l.Printf("Error sending messages: %s\n", err)
	}
}
