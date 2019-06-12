package main

import (
	"flag"
	"fmt"
	"io"
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

// greet the planetosphere
func greet(writer io.Writer) {
	fmt.Fprint(writer, "Hello world!")
}

// handler handles http
func handler(w http.ResponseWriter, r *http.Request) {
	greet(w)
}

func main() {
	flag.Parse()
	fmt.Println("Listening at 127.0.0.1:" + port)
	http.ListenAndServe(":"+port, http.HandlerFunc(handler))
}
