package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"messaging/internal/data"
	"net/http"
)

func (b *Bus) CheckIfBusExists(rw http.ResponseWriter, r *http.Request) {

	bus := data.Bus{}

	read, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(read, &bus)

	if err != nil {
		b.l.Printf("Error unmarshaling request body: %s\n", err)
		http.Error(rw, "Unable to marshal request body", http.StatusBadRequest)
		return
	}

	fmt.Println(bus.Name)
	res, err := data.CheckIfBusExists(r.Context(), data.CheckIfBusExistsInput{
		Client: b.client,
		Bus:    &bus,
	})

	resBody, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	rw.Write(resBody)
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
