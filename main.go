package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
)

func Main(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("index.html")

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	t.Execute(w, req.Host)
}

func main() {
	m := martini.Classic()
	m.Get("/", Main)
	m.Get("/feeds", FetchFeeds)
	m.Post("/feed", AddFeed)

	http.ListenAndServe(":8080", m)
}
