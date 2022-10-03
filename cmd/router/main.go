package main

import (
	"fmt"
	"log"
	routerhandlers "messaging/internal/handlers/router"
	"net/http"
	"os"
)

func main() {

	fmt.Println("SDF")
	l := log.New(os.Stdout, "messaging", log.LstdFlags)

	rh := routerhandlers.NewRouter(l)
	http.HandleFunc("/", rh.ReverseProxy)

	if err := http.ListenAndServe(":1040", nil); err != nil {
		panic(err)
	}

}
