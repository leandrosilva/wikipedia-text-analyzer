package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var analyserURL = getAnalyserURL()
var doneHookURL = getDoneHookURL()

type issueRequest struct {
	URL string `json:"url"`
}

type issueResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

func handleIssue(w http.ResponseWriter, r *http.Request) {
	req, err := getIssueRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := issue(req.URL)
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

func getIssueRequest(r *http.Request) (issueRequest, error) {
	var request issueRequest

	url, found := getQueryParam(r, "url")
	if !found {
		return request, fmt.Errorf("Missing 'url' param in query string")
	}
	request.URL = url

	return request, nil
}

func getQueryParam(r *http.Request, key string) (string, bool) {
	for k, vs := range r.URL.Query() {
		if len(vs) > 0 && strings.ToLower(k) == key {
			return vs[0], true
		}
	}
	return "", false
}

func issue(url string) (issueResponse, error) {
	var response issueResponse

	req, err := json.Marshal(map[string]string{
		"client": "oetacli",
		"url":    url,
		"hook":   doneHookURL})
	if err != nil {
		return response, err
	}

	res, err := http.Post(analyserURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, nil
	}

	return response, nil
}

func getAnalyserURL() string {
	host, found := os.LookupEnv("SERVER_URL")
	if !found {
		host = "http://localhost:8080"
	}
	return host + "/analyse"
}

func getDoneHookURL() string {
	host, found := os.LookupEnv("CLIENT_URL")
	if !found {
		host = "http://localhost:9090"
	}
	return host + "/done"
}
