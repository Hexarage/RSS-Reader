package main

import (
	RSSReader "RSS-Reader/internal"
	"log"
)

func main() {
	// some rudimentary logging
	log.Printf("Starting the main function.\n")
	result := RSSReader.Parse(nil)

	if result == nil {
		log.Fatal("RSS Parser returned empty slice of RSSItems")
	}
}
