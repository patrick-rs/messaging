package main

import (
	"log"
	routerhandlers "messaging/internal/handlers/router"
	"net/http"
	"os"
)

func main() {

	l := log.New(os.Stdout, "messaging", log.LstdFlags)

	rh := routerhandlers.NewRouter(l)
	http.HandleFunc("/", rh.ReverseProxy)

	if err := http.ListenAndServe("localhost:1040", nil); err != nil {
		panic(err)
	}

}
