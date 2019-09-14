package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	apiserver "github.com/andreylm/nats-api-server/pkg/api-server"
	component "github.com/andreylm/nats-component"
	nats "github.com/nats-io/go-nats"
)

var (
	showHelp     bool
	showVersion  bool
	serverListen string
	natsServers  string
	systemTopic  string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: api-server [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.StringVar(&serverListen, "listen", "0.0.0.0:9090", "Network host:port to listen to")
	flag.StringVar(&natsServers, "nats", nats.DefaultURL, "Network host:port to listen to")
	flag.StringVar(&systemTopic, "nats-system-topic", "_NATS_SYSTEM_TOPIC", "Main NATS topic for discover and status usage")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATA Rider API Server v%s", apiserver.Version)
		os.Exit(0)
	}
}

func main() {
	log.Printf("Host: %s, Starting NATS API Server version %s", serverListen, apiserver.Version)
	component := component.NewComponent("api-server")
	component.SetupConnectionToNATS(natsServers, systemTopic)

	server := apiserver.Server{
		Component: component,
	}

	if err := server.ListenAndServe(serverListen); err != nil {
		log.Fatal(err)
	}
}
