package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	component "github.com/andreylm/nats-component"
	"github.com/gorilla/mux"
)

func (s *Server) showServiceStatus(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			makeJSONResponse(w, map[string]interface{}{
				"Error": err.Error(),
			})
		}
	}()

	uuid := mux.Vars(r)["uuid"]
	if uuid == "" {
		err = errors.New("Missing service uuid")
		return
	}

	nc := s.Component.NATS()
	stats := component.Stats{}

	msg, err := nc.Request(fmt.Sprintf("_%s.%s.status", s.Component.SystemTopic(), uuid), nil, 3*time.Second)
	if err != nil {
		return
	}

	if err := json.Unmarshal(msg.Data, &stats); err != nil {
		return
	}

	makeJSONResponse(w, map[string]interface{}{
		"Response": stats,
	})
}
