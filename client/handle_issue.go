package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	// AnalyzerURL - computed getAnalyzerURL()
	AnalyzerURL = getAnalyzerURL()

	// DoneHookURL - computed getDoneHookURL()
	DoneHookURL = getDoneHookURL()
)

// IssueRequest contains the URL to be analyzed
type IssueRequest struct {
	URL       string `json:"url"`
	Sentences int    `json:"sentences"`
	Phrases   int    `json:"phrases"`
}

// IssueResponse is the immediate response confirming the issuing request
type IssueResponse struct {
	Status    string `json:"status"`
	TargetURL string `json:"targetURL"`
}

// AnalyzeRequest is the payload we send to issue a text analysis on the server
type AnalyzeRequest struct {
	Client    string `json:"client"`
	TargetURL string `json:"targetURL"`
	HookURL   string `json:"hookURL"`
	Sentences int    `json:"sentences"`
	Phrases   int    `json:"phrases"`
}

func handleIssue(w http.ResponseWriter, r *http.Request) {
	req, err := getIssueRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := issue(req)
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

func getIssueRequest(r *http.Request) (IssueRequest, error) {
	var request IssueRequest

	vars := mux.Vars(r)
	url, found := vars["url"]
	if !found {
		return request, fmt.Errorf("Missing 'url' param in query string")
	}

	request.URL = url
	request.Sentences = 1
	request.Phrases = 1

	ssent, found := vars["sentences"]
	if found {
		isent, err := strconv.Atoi(ssent)
		if err == nil {
			request.Sentences = isent
		}
	}

	sphra, found := vars["phrases"]
	if found {
		iphra, err := strconv.Atoi(sphra)
		if err == nil {
			request.Phrases = iphra
		}
	}

	return request, nil
}

func issue(request IssueRequest) (IssueResponse, error) {
	var response IssueResponse

	req, err := json.Marshal(AnalyzeRequest{
		Client:    "mehcli",
		TargetURL: request.URL,
		HookURL:   DoneHookURL,
		Sentences: request.Sentences,
		Phrases:   request.Phrases})
	if err != nil {
		return response, err
	}

	res, err := http.Post(AnalyzerURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	if res.StatusCode != 200 {
		return response, fmt.Errorf("Failed to issue article analysis due to server error - %s", body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, nil
	}

	return response, nil
}

func getAnalyzerURL() string {
	host, found := os.LookupEnv("SERVER_URL")
	if !found {
		host = "http://localhost:8080"
	}
	return host + "/analyze"
}

func getDoneHookURL() string {
	host, found := os.LookupEnv("CLIENT_URL")
	if !found {
		host = "http://localhost:9090"
	}
	return host + "/done"
}
