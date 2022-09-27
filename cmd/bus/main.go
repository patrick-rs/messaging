package main

import (
	"context"
	"log"
	shareddata "messaging/internal/data/shared"
	bushandlers "messaging/internal/handlers/bus"
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
	sm := mux.NewRouter()

	bh := bushandlers.NewBus(l, client, validator)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", bh.GetBus)

	s := mhttp.NewHTTPServer(sm, "9091")

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
