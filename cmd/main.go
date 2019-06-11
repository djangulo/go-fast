package main

import (
	"fmt"
	"io"
	"net/http"
)

func greet(writer io.Writer) {
	fmt.Fprint(writer, "Hello world!")
}

func handler(w http.ResponseWriter, r *http.Request) {
	greet(w)
}

func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(handler))
}
