package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type analyseInput struct {
	Client string
	URL    string
	Hook   string
}

type analyseOutput struct {
	Topics  []string
	Summary []string
}

type doneHookResponse struct {
}

func issueAnalysis(input analyseInput) (bool, error) {
	go pullHook(input.Hook)

	return true, nil
}

func analyse(input analyseInput) (analyseOutput, error) {
	return analyseOutput{}, nil
}

func pullHook(url string) (doneHookResponse, error) {
	log.Println("<hook>", url)
	var response doneHookResponse

	req, err := json.Marshal(map[string]string{
		"client": "oetacli",
		"blah":   "yes"})
	if err != nil {
		return response, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(req))
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
