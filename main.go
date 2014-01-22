package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func main() {
	m := martini.Classic()
	m.Use(martini.Static("static"))
	m.Get("/api/feeds", FetchFeeds)
	m.Post("/api/feed", AddFeed)

	http.ListenAndServe(":"+os.Getenv("PORT"), m)
}
