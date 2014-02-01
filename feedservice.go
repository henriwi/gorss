package main

import (
	"github.com/SlyMarbo/rss"
	"errors"
	"fmt"
)

var feed1 = &rss.Feed{UpdateURL: "http://www.vg.no/rss/create.php?categories=20&keywords=&limit=10"}

var feeds = []*rss.Feed{feed1}

func Add(feed *rss.Feed) error {
	if (feedExists(feed)) {
		return errors.New("Feed with UpdateURL already exists")
	}
	feeds = append(feeds, feed)
	return nil
}

func feedExists(feed *rss.Feed) bool {
	for _, v := range feeds {
		if (v.UpdateURL == feed.UpdateURL) {
			return true
		}
	}
	return false
}

func GetAll() []*rss.Feed {
	return feeds
}

func MarkItemUnread(feedIndex int, itemIndex int) {
	feed := feeds[feedIndex]
	item := feed.Items[itemIndex]
	fmt.Printf("Markin %s as read", item)
	item.Read = true
}