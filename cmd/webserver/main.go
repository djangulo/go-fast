package main

import (
	"flag"
	"fmt"
	"github.com/djangulo/go-fast"
	"github.com/djangulo/go-fast/config"
	"log"
	"net/http"
)

var port string

func init() {
	const (
		defaultPort = "9000"
		portUsage   = "port to serve the app on, default '" + defaultPort + "'"
	)
	flag.StringVar(&port, "port", defaultPort, portUsage)
	flag.StringVar(&port, "p", defaultPort, portUsage+" (shorthand)")
}

func main() {
	flag.Parse()
	fmt.Println("Listening at port " + port)
	store, removeStore := poker.NewPostgreSQLPlayerStore(
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabaseName,
		config.DatabasePassword,
	)
	defer removeStore()
	server, err := poker.NewPlayerServer(store)
	if err != nil {
		log.Fatalf("could not create server %v", err)
	}
	if err = http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
