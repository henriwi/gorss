package main

import (
	"github.com/gorilla/mux"
	"github.com/henriwi/gorss/hub"
	"net/http"
)

func main() {
	hub.FetchFeeds()
	r := mux.NewRouter()
	r.HandleFunc("/feeds", hub.FeedsHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
