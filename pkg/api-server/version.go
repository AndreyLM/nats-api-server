package apiserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) showVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Version string `json:"version"`
	}{Version: Version})
	w.Write(response)
}
