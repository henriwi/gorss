package main

import (
	"github.com/codegangsta/martini"
	"net/http"
	"os"
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("static"))
	m.Get("/api/feed", FetchFeeds)
	m.Post("/api/feed", AddFeed)
	m.Delete("/api/feed", DeleteFeed)
	m.Post("/api/feed/read", MarkUnread)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}
