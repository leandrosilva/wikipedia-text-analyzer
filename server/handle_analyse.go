package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// AnalyseRequest is the payload clients send to issue a text analysis
type AnalyseRequest struct {
	Client    string `json:"client"`
	TargetURL string `json:"targetURL"`
	HookURL   string `json:"hookURL"`
	Sentences string `json:"sentences"`
}

// AnalyseResponse is the immediate response we give to clients issuing analysis
type AnalyseResponse struct {
	Status    string `json:"status"`
	TargetURL string `json:"targetURL"`
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

	err = issueAnalysis(AnalyseInput{
		Client:    req.Client,
		TargetURL: req.TargetURL,
		HookURL:   req.HookURL,
		Sentences: getSentences(req.Sentences)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(AnalyseResponse{
		Status:    "issued",
		TargetURL: req.TargetURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getRequest(r *http.Request) (AnalyseRequest, error) {
	var request AnalyseRequest

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(body, &request)

	return request, err
}

func getSentences(s string) int {
	sentences, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return sentences
}
