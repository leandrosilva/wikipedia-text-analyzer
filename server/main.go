package main

import (
	"log"
	"net/http"
	"os"
)

func getPort() string {
	port, found := os.LookupEnv("SERVER_HTTP_PORT")
	if !found {
		port = "8080"
	}
	return port
}

func main() {
	http.HandleFunc("/analyse", handleAnalyse)
	port := getPort()

	log.Println("Starting server at :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
