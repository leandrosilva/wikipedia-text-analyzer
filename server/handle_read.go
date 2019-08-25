package main

import (
	"net/http"
	"net/url"
	"strings"
)

var (
	// ReadURL - computed getReadURL()
	ReadURL = getReadURL()
)

func handleRead(w http.ResponseWriter, r *http.Request) {
	path := splitPath(r.URL)
	articleKey := path[0]

	if articleKey == "" {
		http.Error(w, "Missing article key", http.StatusBadRequest)
		return
	}

	if len(path) > 1 && path[1] == "raw" {
		handleReadRaw(w, r, articleKey)
		return
	}

	handleReadProcessed(w, r, articleKey)
}

func handleReadProcessed(w http.ResponseWriter, r *http.Request, articleKey string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("processed-"))
	w.Write([]byte(articleKey))
}

func handleReadRaw(w http.ResponseWriter, r *http.Request, articleKey string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("raw-"))
	w.Write([]byte(articleKey))
}

func splitPath(url *url.URL) []string {
	return strings.Split(url.Path[len("/read/"):], "/")
}

func getReadURL() string {
	return server.URL + "/read"
}
