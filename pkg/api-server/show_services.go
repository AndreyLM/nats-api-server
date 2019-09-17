package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	component "github.com/andreylm/nats-component"
	"github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

func (s *Server) showServices(w http.ResponseWriter, r *http.Request) {
	count, err := getConnectedServicesCount(s.NATSMonitoringServer + "/varz")
	if err != nil {
		log.Println(err)
	}

	replyTopic := nuid.Next()
	discoveryChan := make(chan *nats.Msg)
	nc := s.Component.NATS()

	sub, err := nc.Subscribe(replyTopic, func(msg *nats.Msg) {
		discoveryChan <- msg
	})
	if err != nil {
		log.Println(err)
		makeJSONResponse(w, map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		cancel()
		close(discoveryChan)
		sub.Unsubscribe()
	}()

	nc.PublishRequest(fmt.Sprintf("_%s.discovery", s.Component.SystemTopic()), replyTopic, nil)

	response := []component.Info{}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case msg := <-discoveryChan:
			info := component.Info{}
			if err = json.Unmarshal(msg.Data, &info); err != nil {
				log.Println(err)
				break
			}
			response = append(response, info)
			if count != 0 && len(response) == count {
				break loop
			}
		}
	}

	makeJSONResponse(w, map[string]interface{}{
		"Services": response,
	})

}

func getConnectedServicesCount(server string) (count int, err error) {
	var data map[string]interface{}
	resp, err := http.Get(server)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	floatCount, ok := data["connections"].(float64)
	if !ok {
		err = errors.New("Error getting connections count")
		return
	}
	count = int(floatCount)
	return
}
