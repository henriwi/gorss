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
	m.Post("/api/feed/:fid/item/:iid/read", MarkUnread)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}
