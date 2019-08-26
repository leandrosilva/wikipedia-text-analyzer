package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	// ReadURL - computed getReadURL()
	ReadURL = getReadURL()
)

func handleRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleKey := vars["articleKey"]

	if articleKey == "" {
		http.Error(w, "Missing article key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("done-"))
	w.Write([]byte(articleKey))
}

func handleReadAtState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleKey := vars["articleKey"]
	state := vars["state"]

	content, err := readArticleContent(articleKey, state)
	if err != nil {
		msg, statusCode := getError(articleKey, err)
		http.Error(w, msg, statusCode)
		return
	}

	switch state {
	case "raw", "summarized":
		w.Header().Set("Content-Type", "text/plain")
	case "ranked":
		w.Header().Set("Content-Type", "application/json")
	}

	w.Write(content)
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
