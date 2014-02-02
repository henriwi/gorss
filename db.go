package main

import (
	"github.com/SlyMarbo/rss"
	"errors"
	"fmt"
)

var feed1 = &rss.Feed{UpdateURL: "http://www.vg.no/rss/create.php?categories=20&keywords=&limit=10"}

var feeds = []*rss.Feed{feed1}

type DB struct {}

func (db DB) Add(feed *rss.Feed) error {
	if (feedExists(feed)) {
		return errors.New("Feed already exists")
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

func (db DB) GetAll() []*rss.Feed {
	return feeds
}

func (db DB) MarkItemUnread(feedIndex int, itemIndex int) {
	feed := feeds[feedIndex]
	item := feed.Items[itemIndex]
	fmt.Printf("Markin %s as read\n", item)
	item.Read = true
}