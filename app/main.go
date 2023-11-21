package main

import (
	RSSReader "RSS-Reader/internal"
	"log"
	"os"
	"time"
)

func main() {
	// some rudimentary logging
	log.Printf("Starting the main function.\n")

	args := os.Args[1:]

	if len(args) == 0 {
		server := NewAPIServer(":3000")
		server.Run()
	}

	config := parseArgs(args)

	/*
		alternatively I can just start a background service which reads a list of rss links and posts a toast to notify when there is a new item
	*/
	// testing out stuff
	log.Printf("Starting rss reader as background service.\n")
	// load in list of links, either from file or from args
	var dummy []string
	var result []RSSReader.RSSItem

	for {
		//TODO: Use the config in order to decide how exactly this runs, something like config.Run() or runDaemon(config)

		newResult := RSSReader.Parse(dummy)

		if slicesAreSame(result, newResult) {
			// call some sort of Toast (can use something like https://github.com/variadico/noti to make it os independent)
			log.Printf("Simulating that the returned list of items is different and we're calling a notification of some sort.\n")
		}
		result = newResult

		time.Sleep(5 * time.Second) // TODO: Consider changing this to a longer duration to reduce network spam
	}

}
