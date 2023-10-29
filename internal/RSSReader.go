package RSSReader

import "fmt"

type RSSItem struct {
}

func (i *RSSItem) Parse(links []string) {
	fmt.Println("Parsing")
}
