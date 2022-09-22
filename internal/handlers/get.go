package handlers

import (
	"encoding/json"
	"io/ioutil"
	"messaging/internal/data"
	"net/http"
)

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

	res, err := data.ReceiveMessages(r.Context(), data.ReceiveMessagesInput{
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
