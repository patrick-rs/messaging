package routerhandlers

import (
	"fmt"
	"log"
	shareddata "messaging/internal/data/shared"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Router struct {
	l *log.Logger
	v *shareddata.Validation
}

func NewRouter(l *log.Logger) *Router {
	return &Router{l: l}
}

func (m *Router) ReverseProxy(rw http.ResponseWriter, r *http.Request) {
	port := ""
	host := ""
	switch r.URL.Path {
	case "/message":
		port = os.Getenv("MESSAGE_PORT")
		host = "message"
	case "/bus":
		port = os.Getenv("BUS_PORT")
		host = "bus"
	}

	fmt.Println("MESSAGE_PORT", host, port)
	target, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))

	if err != nil {
		m.l.Printf("Error parsing url, host: %s, port: %s\n", host, port)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	r.Host = r.URL.Host

	proxy.ServeHTTP(rw, r)
}
