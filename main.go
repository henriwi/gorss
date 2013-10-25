package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

type Item struct {
	Title string `xml:"title"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Link string `xml:"link"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Result struct {
	Channel Channel `xml:"channel"`
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadFile("resources/feed.rss")

  if (err != nil) {
  	log.Fatal(err)
  }

  var r Result
	_ = xml.Unmarshal(data, &r)

	jsonValue, _ := json.Marshal(r)

	fmt.Fprintf(writer, string(jsonValue))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}