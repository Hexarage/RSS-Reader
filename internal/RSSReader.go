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
	Title       string `xml:"title,omitempty"`
	Source      string `xml:"source,omitempty"`
	SourceURL   string `xml:"sourceurl,omitempty"` // haven't encountered this particular tag, TODO: See what it may be referring to
	Link        string `xml:"link"`
	PublishDate string `xml:"pubDate"` // running into some weird parsing error for Time.time, so keeping it as string for now
	Description string `xml:"description"`
}
type channel struct {
	Title string    `xml:"title,omitempty"`
	Image string    `xml:"image,omitempty"`
	Items []RSSItem `xml:"item"`
}

type rSSFeed struct { // figure out how to unmarshal into this?
	Channel channel `xml:"channel"`
} // using https://en.wikipedia.org/wiki/RSS#Example as a template of sorts

func checkF(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// add a context so I can cancel mid way?
func parseRSSLink(link *url.URL, ch chan<- []RSSItem, wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO: Fetch the data from the link
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(link.String())
	checkF(err)

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	checkF(err)
	//log.Printf("Successfully read the data from the link %v\n %v", link, string(data))
	// TODO: do some parsing and stuff

	items := parseFeed(data)
	if items != nil {
		ch <- items
	}
	// pass the parsed item to the channel if valid, otherwise just be done with this go routine

}
func parseFeed(data []byte) []RSSItem {

	/*
		need to traverse the feed until I get to items
		after which parse each item and pass them on
	*/

	var f rSSFeed
	//var test map[string]interface{}
	err := xml.Unmarshal(data, &f) // RSS is just an xml feed
	checkF(err)

	return f.Channel.Items
}

func Parse(links []string) []RSSItem {
	if links == nil {
		return nil
	}

	waitGroup := sync.WaitGroup{}
	rssChannel := make(chan []RSSItem)

	for _, link := range links {
		// check if link is a valid link
		currentLink, err := url.ParseRequestURI(link)
		checkF(err)

		waitGroup.Add(1)
		go parseRSSLink(currentLink, rssChannel, &waitGroup)
	}

	go func() { // TODO:  really hacky, get rid of this later
		waitGroup.Wait()
		close(rssChannel)
	}()

	var parsedRSSItems []RSSItem

	for item := range rssChannel {
		parsedRSSItems = append(parsedRSSItems, item...)
	}

	return parsedRSSItems
}
