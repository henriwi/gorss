package main

import (
	"encoding/json"
	"fmt"
	"github.com/SlyMarbo/rss"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpResponse struct {
	feed *rss.Feed
	err  error
}

var db = new(DB)

func FetchFeeds(writer http.ResponseWriter, request *http.Request) {
	rss.CacheParsedItemIDs(false)
	feeds, err := db.GetAll()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	if len(feeds) == 0 {
		return
	}

	responses := asyncFetchFeeds(feeds)

	result := []*rss.Feed{}
	for _, r := range responses {
			result = append(result, r.feed)
	}

	jsonResult, _ := json.Marshal(result)
	fmt.Fprintf(writer, string(jsonResult))
}

func asyncFetchFeeds(feeds []*rss.Feed) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for _, feed := range feeds {
		go func(f *rss.Feed) {
			fmt.Printf("Fetching %s\n", f.UpdateURL)
			feed, err := rss.Fetch(f.UpdateURL)

			if err != nil {
				fmt.Printf("Error in response %s. Using old feed.\n", err)
				ch <- &HttpResponse{f, err}
			} else {
				ch <- &HttpResponse{feed, err}
			}

		}(feed)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.feed.UpdateURL)
			responses = append(responses, r)
			if len(responses) == len(feeds) {
				return responses
			}
		case <-time.After(5 * time.Second):
			return responses
		}
	}
}

func AddFeed(writer http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var objmap map[string]string
	json.Unmarshal([]byte(body), &objmap)

	var url = objmap["url"]
	feed, err := rss.Fetch(url)

	if err != nil {
		fmt.Printf("Error fetching feed %s\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.Add(feed)

	if err != nil {
		fmt.Printf("%s\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	jsonResult, _ := json.Marshal(feed)
	fmt.Fprintf(writer, string(jsonResult))
}

func DeleteFeed(writer http.ResponseWriter, req *http.Request, params martini.Params) {
	body, _ := ioutil.ReadAll(req.Body)
	var objmap map[string]string
	json.Unmarshal([]byte(body), &objmap)

	var url = objmap["url"]

	err := db.DeleteFeed(url)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}

func MarkUnread(writer http.ResponseWriter, req *http.Request, params martini.Params) {
	body, _ := ioutil.ReadAll(req.Body)
	var objmap map[string]string
	json.Unmarshal([]byte(body), &objmap)

	var id = objmap["id"]

	err := db.MarkItemUnread(id)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}
