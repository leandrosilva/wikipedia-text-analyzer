package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type doneRequest struct {
	URL string `json:"url"`
}

type doneResponse struct {
	Success bool `json:"success"`
}

func handleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := getDoneRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := done(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getDoneRequest(r *http.Request) (doneRequest, error) {
	log.Println("<hook>")
	var request doneRequest

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(body, &request)

	return request, err
}

func done(url string) (doneResponse, error) {
	var response doneResponse

	return response, nil
}
