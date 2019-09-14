package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nuid"

	"github.com/nats-io/go-nats"
)

func (s *Server) showVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(struct {
		Version string `json:"version"`
	}{Version: Version})
	w.Write(response)
}

func (s *Server) showServices(w http.ResponseWriter, r *http.Request) {
	replyTopic := nuid.Next()
	discoveryChan := make(chan *nats.Msg)
	nc := s.Component.NATS()

	sub, err := nc.Subscribe(replyTopic, func(msg *nats.Msg) {
		discoveryChan <- msg
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		cancel()
		close(discoveryChan)
		sub.Unsubscribe()
	}()

	if err != nil {
		log.Println(err)
		return
	}

	nc.PublishRequest("_NATS_RIDER.discovery", replyTopic, nil)

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Fprint(w, "Discovery time finished...")
			break loop
		case msg := <-discoveryChan:
			fmt.Fprintf(w, "ID: "+string(msg.Data)+"\n")
		}
	}

	w.Header().Set("Content-Type", "application/json")

}
