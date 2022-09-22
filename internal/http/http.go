package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewHTTPServer(sm *mux.Router, port string) http.Server {
	return http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
}
