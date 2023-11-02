package RSSReader

import (
	"net/url"
	"sync"
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

// add a context so I can cancel mid way?
func parseRSSLink(link *url.URL, ch chan<- RSSItem, wg *sync.WaitGroup) {
	defer wg.Done()
	var item RSSItem

	// TODO: do some parsing and stuff

	// pass the parsed item to the channel if valid, otherwise just be done with this go routine
	ch <- item

}

func Parse(links []string) []RSSItem {
	if links == nil {
		return nil
	}

	waitGroup := sync.WaitGroup{}
	rssChannel := make(chan RSSItem)

	for _, link := range links {
		// check if link is a valid link
		currentLink, err := url.ParseRequestURI(link)
		if err != nil {
			panic(err)
		}
		// link is valid, add to a list of actual links?
		// go parseRSSLink(currentLink) + some sort of channel?
		waitGroup.Add(1)
		go parseRSSLink(currentLink, rssChannel, &waitGroup)
	}

	go func() { // TODO:  really hacky, get rid of this later
		waitGroup.Wait()
		close(rssChannel)
	}()

	var parsedRSSItems []RSSItem

	for item := range rssChannel {
		parsedRSSItems = append(parsedRSSItems, item)
	}

	return parsedRSSItems
}

/*
	parse the given links, error out if they are invalid?
	parse the feed of all links asynchronously

*/
