package main

import (
	"flag"
	"fmt"
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

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 3
}
func (i *InMemoryPlayerStore) RecordWin(name string) {}

func main() {
	flag.Parse()
	fmt.Println("Listening at 127.0.0.1:" + port)
	server := &PlayerServer{&InMemoryPlayerStore{}}
	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
