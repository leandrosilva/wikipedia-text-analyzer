package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// AnalyzeRequest is the payload clients send to issue a text analysis
type AnalyzeRequest struct {
	Client    string `json:"client"`
	TargetURL string `json:"targetURL"`
	HookURL   string `json:"hookURL"`
	Sentences int    `json:"sentences"`
	Phrases   int    `json:"phrases"`
}

// AnalyzeResponse is the immediate response we give to clients issuing analysis
type AnalyzeResponse struct {
	Status    string `json:"status"`
	TargetURL string `json:"targetURL"`
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := getRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = issueAnalysis(AnalyzeInput{
		Client:    req.Client,
		TargetURL: req.TargetURL,
		HookURL:   req.HookURL,
		Sentences: req.Sentences,
		Phrases:   req.Phrases})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(AnalyzeResponse{
		Status:    "issued",
		TargetURL: req.TargetURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getRequest(r *http.Request) (AnalyzeRequest, error) {
	var request AnalyzeRequest

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(body, &request)

	return request, err
}
