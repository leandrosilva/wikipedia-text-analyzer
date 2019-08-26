package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// ServerInfo packs info about this HTTP server
type ServerInfo struct {
	Protocol string
	Host     string
	Port     string
	URL      string
}

var server = ServerInfo{
	Protocol: getProtocol(),
	Host:     getHost(),
	Port:     getPort(),
	URL:      getURL(),
}

func getProtocol() string {
	prot, found := os.LookupEnv("SERVER_HTTP_PROT")
	if !found {
		return "http"
	}
	return prot
}

func getHost() string {
	host, found := os.LookupEnv("SERVER_HTTP_HOST")
	if !found {
		return "localhost"
	}
	return host
}

func getPort() string {
	port, found := os.LookupEnv("SERVER_HTTP_PORT")
	if !found {
		return "8080"
	}
	return port
}

func getURL() string {
	url, found := os.LookupEnv("SERVER_URL")
	if !found {
		return "http://localhost:8080"
	}
	return url
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/analyze", handleAnalyze).Methods("POST")
	router.HandleFunc("/read/{articleKey}", handleRead)
	router.HandleFunc("/read/{articleKey}/{state}", handleReadAtState)

	log.Println("Starting server at :" + server.Port)
	log.Fatal(http.ListenAndServe(":"+server.Port, router))
}
