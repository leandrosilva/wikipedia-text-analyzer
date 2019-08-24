package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type analyseRequest struct {
	Client string `json:"client"`
	URL    string `json:"url"`
	Hook   string `json:"hook"`
}

type analyseResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

func handleAnalyse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success, err := issueAnalysis(analyseInput{
		Client: req.Client,
		URL:    req.URL,
		Hook:   req.Hook})
	if !success || err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(analyseResponse{
		Status: "issued",
		URL:    req.URL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getRequest(r *http.Request) (analyseRequest, error) {
	var request analyseRequest

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(body, &request)

	return request, err
}
