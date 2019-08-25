package main

import (
	"net/http"
)

var (
	// ReadURL - computed getReadURL()
	ReadURL = getReadURL()
)

func handleRead(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/read/"):]

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(key))
}

func getReadURL() string {
	return server.URL + "/read"
}
