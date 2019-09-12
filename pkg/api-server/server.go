package apiserver

import component "github.com/andreylm/nats-component"

// Version - api server version
var Version = "v1.0.0"

// Server - api server
type Server struct {
	Component *component.Component
}

// ListenAndServe - listen and serve
func (s *Server) ListenAndServe(host string) error {
	return nil
}
