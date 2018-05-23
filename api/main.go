package main

import (
	"log"

	"github.com/RaniSputnik/ok/api/app"
)

func main() {
	addr := ":8080"
	server := app.New(app.Config{Addr: addr})
	log.Printf("Now accepting connections at http://localhost%s", addr)
	log.Fatal(server.ListenAndServe())
}
