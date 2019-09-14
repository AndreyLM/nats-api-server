package apiserver

import (
	"net/http"
	"time"

	component "github.com/andreylm/nats-component"
	"github.com/gorilla/mux"
)

// Version - api server version
var Version = "v1.0.0"

// Server - api server
type Server struct {
	Component *component.Component
}

// ListenAndServe - listen and serve
func (s *Server) ListenAndServe(host string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", s.showVersion)
	r.HandleFunc("/services", s.showServices)

	srv := &http.Server{
		Handler:      r,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}
