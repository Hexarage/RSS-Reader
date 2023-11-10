package main

import (
	"log"
)

func main() {
	// some rudimentary logging
	log.Printf("Starting the main function.\n")

	server := NewAPIServer(":3000")
	server.Run()
}
