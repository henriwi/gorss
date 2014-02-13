package main

import (
	"errors"
	"fmt"
	"github.com/SlyMarbo/rss"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"strings"
)

type DB struct{}

var (
	mgoSession   *mgo.Session
	dbName string
)

func getSession() (*mgo.Session, error) {
	if mgoSession == nil {
		url := os.Getenv("MONGOHQ_URL")

		if url == "" {
			fmt.Println("Connection url is empty\n")
			return nil, errors.New("Connection url is empty")
		}

		slashIndex := strings.LastIndex(url, "/")
		dbName = url[slashIndex + 1:len(url)]

		var err error
		mgoSession, err = mgo.Dial(url)
		if err != nil {
			fmt.Printf("Error when connecting to db %s\n", err)
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
}

func (db DB) Add(feed *rss.Feed) error {
	sess, err := getSession()

	if err != nil {
		return err
	}

	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")

	var existingFeed *rss.Feed
	collection.Find(bson.M{"updateurl": feed.UpdateURL}).One(&existingFeed)

	if existingFeed != nil {
		return errors.New("Feed already exists")
	}

	err = collection.Insert(feed)
	if err != nil {
		fmt.Printf("Can't insert feed: %v\n", err)
		return err
	}

	return nil
}

func (db DB) GetAll() ([]*rss.Feed, error) {
	sess, err := getSession()

	if err != nil {
		return nil, err
	}

	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")

	var feeds []*rss.Feed
	collection.Find(nil).All(&feeds)
	return feeds, nil
}

func (db DB) DeleteFeed(updateURL string)  error {
	sess, err := getSession()

	if err != nil {
		return err
	}

	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")

	err = collection.Remove(bson.M{"updateurl": updateURL})

	if err != nil {
		fmt.Printf("Error deleting feed: %v\n", err)
		return err
	}

	return nil
}

func (db DB) MarkItemUnread(id string) error {
	sess, err := getSession()

	if err != nil {
		return err
	}

	defer sess.Close()

	collection := sess.DB(dbName).C("feeds")
	err = collection.Update(bson.M{"items.id": id}, bson.M{"$set": bson.M{"items.$.read": true}})

	if err != nil {
		fmt.Printf("Error marking feed as read: %v\n", err)
		return err
	}

	return nil
}
