package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/SlyMarbo/rss"
	"io/ioutil"
)

type HttpResponse struct {
	feed *rss.Feed
	err  error
}

var feed1 = &rss.Feed{UpdateURL: "http://www.vg.no/rss/create.php?categories=20&keywords=&limit=10"}

var feeds = []*rss.Feed{feed1}

// "http://www.aftenposten.no/rss/?kat=nyheter_iriks" : nil,

func FetchFeeds(writer http.ResponseWriter, request *http.Request) {
	responses := asyncFetchFeeds()

	result := []*rss.Feed{}
	for _, r := range responses {
		if r.err != nil {
			fmt.Printf("Error in response %s\n", r.err)
		} else {
			result = append(result, r.feed)
		}
	}

	jsonResult, _ := json.Marshal(result)
	fmt.Fprintf(writer, string(jsonResult))
}

func asyncFetchFeeds() []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for _, feed := range feeds {
		fmt.Printf("Fetching %s\n", feed)
			go func(feed *rss.Feed) {
				err := feed.Update()
				ch <- &HttpResponse{feed, err}
			}(feed)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.feed)
			responses = append(responses, r)
			if len(responses) == len(feeds) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

func AddFeed(writer http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var objmap map[string]string
	json.Unmarshal([]byte(body), &objmap)

	var newUrl = objmap["url"]
	newFeed := &rss.Feed{UpdateURL: newUrl}
	feeds = append(feeds, newFeed)

 	writer.WriteHeader(http.StatusCreated)
}
