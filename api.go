package main

import (
	"encoding/json"
	"fmt"
	"github.com/SlyMarbo/rss"
	"net/http"
	"time"
)

type HttpResponse struct {
	url  string
	feed *rss.Feed
	err  error
}

var urls = map[string]*rss.Feed{
	"http://www.aftenposten.no/rss/?kat=nyheter_iriks" : nil,
	"http://www.vg.no/rss/create.php?categories=20&keywords=&limit=10": nil,
}

func FetchFeeds(writer http.ResponseWriter, request *http.Request) {
	responses := asyncFetchFeeds()

	feeds := []*rss.Feed{}
	for _, r := range responses {
		if r.err != nil {
			fmt.Printf("Error in response %s\n", r.err)
		} else {
			feeds = append(feeds, r.feed)
		}
	}

	jsonResult, _ := json.Marshal(feeds)
	fmt.Fprintf(writer, string(jsonResult))
}

func asyncFetchFeeds() []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for url, feed := range urls {
		if feed != nil {
			fmt.Printf("Updating %s\n", url)
			go func(url string, feed *rss.Feed) {
				err := feed.Update()
				ch <- &HttpResponse{url, feed, err}
			}(url, feed)
		} else {
			fmt.Printf("Fetching new %s\n", url)
			go func(url string) {
				fmt.Printf("Fetching %s \n", url)
				feed, err := rss.Fetch(url)
				urls[url] = feed
				ch <- &HttpResponse{url, feed, err}
			}(url)
		}
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.url)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

func AddFeed(req *http.Request) int {
	url := req.FormValue("url")
	fmt.Printf("Adding %s\n", url)
	urls[url] = nil
	return http.StatusOK
}
