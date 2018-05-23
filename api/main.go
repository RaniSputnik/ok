package main

import (
	"log"
	"net/http"

	"github.com/RaniSputnik/ok/api/handle"
)

func main() {
	// TODO create server
	addr := ":8080"
	log.Printf("Now accepting connections at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, handle.New()))
}
