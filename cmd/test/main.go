package main

import (
	"log"

	apiserver "github.com/andreylm/nats-api-server/pkg/api-server"
)

func main() {
	log.Println(apiserver.Version)
}
