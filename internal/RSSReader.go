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
	Title       string `xml:"title,omitempty" json:"title"`
	Source      string `xml:"source,omitempty" json:"source"`
	SourceURL   string `xml:"sourceurl,omitempty" json:"source_url"` // haven't encountered this particular tag, TODO: See what it may be referring to
	Link        string `xml:"link,omitempty" json:"link"`
	PublishDate string `xml:"pubDate,omitempty" json:"publish_date"` // running into some weird parsing error for Time.time, so keeping it as string for now
	Description string `xml:"description,omitempty" json:"description"`
}

type channel struct {
	Title string    `xml:"title,omitempty"`
	Image string    `xml:"image,omitempty"`
	Items []RSSItem `xml:"item"`
}

type rSSFeed struct { // figure out how to unmarshal into this?
	Channel channel `xml:"channel"`
}

func parseRSSLink(link *url.URL, ch chan<- []RSSItem, wg *sync.WaitGroup) {
	defer wg.Done()

	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(link.String())
	if err != nil {
		log.Printf("Encountered an error: \"%v\" while getting data from link: %v", err, link.String())
		return
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Encountered an error: \"%v\" while reading data from link: %v", err, link.String())
		return
	}

	items := parseFeed(data)
	if items != nil {
		for _, i := range items {
			i.SourceURL = link.String()
		}
		ch <- items
	}
}

func parseFeed(data []byte) []RSSItem {
	var f rSSFeed

	err := xml.Unmarshal(data, &f)
	if err != nil {
		log.Printf("Encountered an error: \"%v\" while unmarshaling xml data.", err)
		return nil
	}

	return f.Channel.Items
}

func Parse(links []string) []RSSItem {
	if links == nil {
		log.Printf("Parser was passed an empty slice of links")
		return nil
	}

	waitGroup := sync.WaitGroup{}
	rssChannel := make(chan []RSSItem)

	for _, link := range links {
		// check if link is a valid link
		currentLink, err := url.ParseRequestURI(link)
		if err != nil {
			log.Printf("Encountered an error: \"%v\" while verifying link validity for link:\"%v\"", err, link)
			return nil
		}

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
