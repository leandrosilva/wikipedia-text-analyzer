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
		handleReadOfRaw(w, r, articleKey)
		return
	}

	handleReadOfDone(w, r, articleKey)
}

func handleReadOfDone(w http.ResponseWriter, r *http.Request, articleKey string) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("done-"))
	w.Write([]byte(articleKey))
}

func handleReadOfRaw(w http.ResponseWriter, r *http.Request, articleKey string) {
	content, err := readWikipediaRawArticle(articleKey)
	if err != nil {
		msg, statusCode := getError(articleKey, err)
		http.Error(w, msg, statusCode)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func splitPath(url *url.URL) []string {
	return strings.Split(url.Path[len("/read/"):], "/")
}

func getError(articleKey string, err error) (string, int) {
	msg := err.Error()
	if strings.Contains(msg, "The system cannot find the file specified") {
		return "Key " + articleKey + " is unknown.", http.StatusNotFound
	}
	return msg, http.StatusInternalServerError
}

func getReadURL() string {
	return server.URL + "/read"
}
