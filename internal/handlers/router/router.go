package routerhandlers

import (
	"fmt"
	"log"
	shareddata "messaging/internal/data/shared"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Router struct {
	l *log.Logger
	v *shareddata.Validation
}

func NewRouter(l *log.Logger) *Router {
	return &Router{l: l}
}

func (m *Router) ReverseProxy(rw http.ResponseWriter, r *http.Request) {
	port := 0
	switch r.URL.Path {
	case "/message":
		port = 1041
	case "/bus":
		port = 1042
	}

	target, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))

	if err != nil {
		m.l.Printf("Error parsing url, port: %d\n", port)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	r.Host = r.URL.Host

	proxy.ServeHTTP(rw, r)
}
