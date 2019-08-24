package main

import (
	"log"
	"net/http"
	"os"
)

func getPort() string {
	port, found := os.LookupEnv("CLIENT_HTTP_PORT")
	if !found {
		port = "9090"
	}
	return port
}

func main() {
	http.HandleFunc("/issue", handleIssue)
	http.HandleFunc("/done", handleDone)
	port := getPort()

	log.Println("Starting server at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
