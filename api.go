package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/SlyMarbo/rss"
	"io/ioutil"
	"github.com/codegangsta/martini"
	"strconv"
)

type HttpResponse struct {
	feed *rss.Feed
	err  error
}

var db = new(DB)

func FetchFeeds(writer http.ResponseWriter, request *http.Request) {
	feeds := db.GetAll()

	if (len(feeds) == 0) {
		return
	}

	responses := asyncFetchFeeds(feeds)

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

func asyncFetchFeeds(feeds []*rss.Feed) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for _, feed := range feeds {
		fmt.Printf("Fetching %s\n", feed.UpdateURL)
			go func(feed *rss.Feed) {
				err := feed.Update()
				time.Sleep(500 * time.Millisecond)
				ch <- &HttpResponse{feed, err}
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

	var url = objmap["url"]
	feed, err := rss.Fetch(url)

	if (err != nil) {
		fmt.Printf("Error fetching feed %s", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.Add(feed)

	if (err != nil) {
		fmt.Printf("%s\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

 	writer.WriteHeader(http.StatusCreated)
	jsonResult, _ := json.Marshal(feed)
	fmt.Fprintf(writer, string(jsonResult))
}

func DeleteFeed(writer http.ResponseWriter, req *http.Request, params martini.Params) {
	feedIndex, _ := strconv.Atoi(params["id"])
	db.DeleteFeed(feedIndex)
	writer.WriteHeader(http.StatusOK)
}

func MarkUnread(writer http.ResponseWriter, req *http.Request, params martini.Params) {
	feedIndex, _ := strconv.Atoi(params["_0"]);
	itemIndex, _ := strconv.Atoi(params["_1"]);

	db.MarkItemUnread(feedIndex, itemIndex)
	writer.WriteHeader(http.StatusOK)
}