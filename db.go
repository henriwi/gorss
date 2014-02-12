package main

import (
	"errors"
	"fmt"
	"github.com/SlyMarbo/rss"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type DB struct{}

var (
	mgoSession   *mgo.Session
	dbName = "simplyrss"
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		url := os.Getenv("MONGOHQ_URL")

		if url == "" {
			fmt.Println("Connection url is empty\n")
		}

		var err error
		mgoSession, err = mgo.Dial(url)
		if err != nil {
			fmt.Printf("Error when connecting to db %s\n", err)
		}
	}
	return mgoSession.Clone()
}

func (db DB) Add(feed *rss.Feed) error {
	sess := getSession()
	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")

	var existingFeed *rss.Feed
	collection.Find(bson.M{"updateurl": feed.UpdateURL}).One(&existingFeed)

	if existingFeed != nil {
		return errors.New("Feed already exists")
	}

	err := collection.Insert(feed)
	if err != nil {
		fmt.Printf("Can't insert feed: %v\n", err)
		return err
	}

	return nil
}

func (db DB) GetAll() []*rss.Feed {
	sess := getSession()
	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")

	var feeds []*rss.Feed
	collection.Find(nil).All(&feeds)
	return feeds
}

func (db DB) DeleteFeed(updateURL string) {
	sess := getSession()
	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")
	
	_ = collection.Remove(bson.M{"updateurl": updateURL})
}

func (db DB) MarkItemUnread(id string) {
	sess := getSession()
	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")
	_ = collection.Update(bson.M{"items.id": id}, bson.M{"$set": bson.M{"items.$.read": true}})
}
