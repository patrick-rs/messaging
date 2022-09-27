package main

import (
	"context"
	"log"
	shareddata "messaging/internal/data/shared"
	messagehandlers "messaging/internal/handlers/message"
	mhttp "messaging/internal/http"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	client, err := shareddata.NewMongoDBClient()
	if err != nil {
		l.Printf("Error creating new mongodb client: %s", err)
		return
	}

	validator := shareddata.NewValidation()

	mh := messagehandlers.NewMessages(l, client, validator)

	sm := mux.NewRouter()

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", mh.PostMessage)
	postRouter.Use(mh.MiddlewareMessagesValidation)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", mh.GetMessage)

	s := mhttp.NewHTTPServer(sm, "9090")

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Received termination, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), time.Second*30)
	s.Shutdown(tc)

}
