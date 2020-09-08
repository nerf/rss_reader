package rss_reader

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nerf/rss_reader/utils"
)

// RssItem struct that contains rss item data
type RssItem struct {
	Title       string    `json:"title"`
	Source      string    `json:"source"`
	SourceURL   string    `json:"source_url"`
	Link        string    `json:"link"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`
}

type ping struct{}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	XMLName xml.Name  `xml:"channel"`
	Source  string    `xml:"title"`
	Items   []rssItem `xml:"item"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PublishDate string   `xml:"pubDate"`
}

// Parse a list of urls and return found items combined
// Silently ignore broken url's
func Parse(urls []string) (items []RssItem, err error) {
	numberOfUrls := len(urls)

	if numberOfUrls == 0 {
		return items, errors.New("Urls list is empty")
	}

	itemsCh := make(chan RssItem)
	workersQuitCh := make(chan ping)
	workersDoneCount := 0

	for _, url := range urls {
		go fetchItemsFromURL(url, itemsCh, workersQuitCh)
	}

	for {
		select {
		case item := <-itemsCh:
			items = append(items, item)
		case <-workersQuitCh:
			workersDoneCount++
		}

		if workersDoneCount == numberOfUrls {
			close(itemsCh)
			close(workersQuitCh)
			break
		}
	}

	return items, nil
}

func fetchItemsFromURL(url string, itemsCh chan<- RssItem, quitCh chan<- ping) {
	defer func() {
		quitCh <- ping{}
	}()

	response, err := http.Get(url)
	if err != nil || response.StatusCode > 299 || response.StatusCode < 200 {
		// Silently ignore
		return
	}
	defer response.Body.Close()

	var rssFeed rssFeed
	body, _ := ioutil.ReadAll(response.Body)
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		// Silently ignore errors
		return
	}

	source := rssFeed.Channel.Source
	for _, i := range rssFeed.Channel.Items {
		date, err := utils.ParseDate(i.PublishDate)

		if err != nil {
			// We don't want entries without date
			return
		}

		itemsCh <- RssItem{
			Title:       i.Title,
			Description: i.Description,
			Link:        i.Link,
			PublishDate: date,
			Source:      source,
			SourceURL:   url,
		}
	}
}
