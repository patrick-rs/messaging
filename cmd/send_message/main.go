package main

import (
	"context"
	"log"
	"messaging/internal/data"
	"messaging/internal/handlers"
	mhttp "messaging/internal/http"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type Message struct {
	Queue   string `bson:"queue"`
	Message string `bson:"message"`
}

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	client, err := data.NewMongoDBClient()
	if err != nil {
		l.Printf("Error creating new mongodb client: %s", err)
		return
	}

	validator := data.NewValidation()

	mh := handlers.NewMessages(l, client, validator)

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
