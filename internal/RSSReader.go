package RSSReader

import (
	"fmt"
	"time"
)

type RSSItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	Description string
}

func Parse(links []string) []RSSItem {
	fmt.Println("Parsing")

	return nil
}

/*
	parse the given links, error out if they are invalid?
	parse the feed of all links asynchronously

*/
