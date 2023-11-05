package RSSReader

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
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

type RSSFeed struct { // figure out how to unmarshal into this?
	Channel struct {
		Title       string `xml:"title"`
		Source      string `xml:"source"`
		SourceURL   string `xml:""`
		Link        string
		PublishDate time.Time `xml:"pubDate"`
		Item        struct {
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
} // using https://en.wikipedia.org/wiki/RSS#Example as a template of sorts

func checkF(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// add a context so I can cancel mid way?
func parseRSSLink(link *url.URL, ch chan<- RSSItem, wg *sync.WaitGroup) {
	defer wg.Done()
	var item RSSItem

	// TODO: Fetch the data from the link
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(link.String())
	checkF(err)

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	checkF(err)
	log.Printf("Successfully read the data from the link %v\n %v", link, string(data))
	// TODO: do some parsing and stuff

	var feed RSSFeed

	err = xml.Unmarshal(data, &feed) // RSS is just an xml feed
	checkF(err)
	//log.Println(RSSFeed.Channel.Title)
	parseFeed(data)

	// pass the parsed item to the channel if valid, otherwise just be done with this go routine
	ch <- item

}
func parseFeed(data []byte) RSSFeed {

	return RSSFeed{}
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
		checkF(err)

		// link is valid, add to a list of actual links?
		// go parseRSSLink(currentLink) + some sort of channel?
		waitGroup.Add(1)
		go parseRSSLink(currentLink, rssChannel, &waitGroup)
	}

	go func() { // TODO:  really hacky, get rid of this later, maybe just make a channel that is as big as the ammount of links passed?
		waitGroup.Wait()
		close(rssChannel)
	}()

	var parsedRSSItems []RSSItem

	for item := range rssChannel {
		parsedRSSItems = append(parsedRSSItems, item)
	}

	return parsedRSSItems
}
