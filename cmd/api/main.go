package main

import (
	"log"
)

func main() {

	server := NewApiServer(":8080", nil)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
