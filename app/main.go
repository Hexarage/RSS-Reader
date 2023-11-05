package main

import (
	RSSReader "RSS-Reader/internal"
	"log"
)

func main() {
	// some rudimentary logging
	log.Printf("Starting the main function.\n")
	var links []string
	links = append(links, "https://rss.com/blog/category/press-releases/feed/")
	links = append(links, "https://blog.jetbrains.com/go/feed")

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	server.Run()

	result := RSSReader.Parse(links)

	if result == nil {
		log.Fatal("RSS Parser returned empty slice of RSSItems")
	}

	log.Println("Parsed ", len(result), " items")
}
