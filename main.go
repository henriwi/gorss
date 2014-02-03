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
	m.Delete("/api/feed/:id", DeleteFeed)
	m.Post("/api/feed/:id/item/:id/read", MarkUnread)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}
