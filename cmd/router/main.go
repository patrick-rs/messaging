package main

import (
	"fmt"
	"log"
	routerhandlers "messaging/internal/handlers/router"
	"net/http"
	"os"
)

func main() {

	l := log.New(os.Stdout, "messaging", log.LstdFlags)

	rh := routerhandlers.NewRouter(l)
	http.HandleFunc("/", rh.ReverseProxy)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), nil); err != nil {
		panic(err)
	}

}
