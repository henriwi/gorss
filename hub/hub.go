package hub

import (
	"encoding/json"
	"fmt"
	"github.com/SlyMarbo/rss"
	"net/http"
	"time"
)

var pending []*rss.Item

func FeedsHandler(writer http.ResponseWriter, request *http.Request) {
	jsonResult, _ := json.Marshal(pending)
	fmt.Fprintf(writer, string(jsonResult))
}

func FetchFeeds() {
	urls := []string{"http://localhost:8081", "http://www.vg.no/rss/create.php?categories=20&keywords=&limit=10"}

	for _, url := range urls {
		go fetchFeed(url)
	}
}

func fetchFeed(url string) {
	for {
		feed, err := rss.Fetch(url)

		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		fmt.Printf("There are %d items in %s\n", len(feed.Items), url)

		for _, item := range feed.Items {
			pending = append(pending, item)
		}

		<-time.After(time.Duration(10 * time.Second))
	}
}
