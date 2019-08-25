package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// AnalyseInput is the payload to issue a text analysis
type AnalyseInput struct {
	Client    string
	TargetURL string
	HookURL   string
	Sentences int
	Phrases   int
}

// AnalyseOutput is the final response of a performed text analysis
type AnalyseOutput struct {
	TargetURL string
	ResultURL string
	Topics    []string
	Summary   []string
}

// DoneHookResponse is what clients respond when we pull their hook (ouch!)
type DoneHookResponse struct {
	Acknowledge bool `json:"acknowledge"`
}

func issueAnalysis(input AnalyseInput) error {
	log.Println("Issuing analysis:", input.TargetURL)
	go analyse(input)

	return nil
}

func analyse(input AnalyseInput) (AnalyseOutput, error) {
	log.Println("Starting analysis:", input.TargetURL)

	// Get data
	articleKey := getArticleKey(input.TargetURL)
	rawFilePath, err := downloadWikipediaArticle(input.TargetURL, articleKey)
	if err != nil {
		return AnalyseOutput{}, err
	}

	// Summarize it
	_, err = summarizeArticle(articleKey, rawFilePath, input.Sentences)
	if err != nil {
		return AnalyseOutput{}, err
	}

	// Rank k phrases
	_, err = rankArticlePhrases(articleKey, rawFilePath, input.Phrases)
	if err != nil {
		return AnalyseOutput{}, err
	}

	// Clean text
	_, err = cleanArticle(articleKey, rawFilePath)
	if err != nil {
		return AnalyseOutput{}, err
	}

	log.Println("Finished analysis:", input.TargetURL)

	output := AnalyseOutput{
		TargetURL: input.TargetURL,
		ResultURL: getResultURL(articleKey)}

	res, err := pullHook(input.HookURL, output)
	if err != nil {
		return output, err
	}
	log.Println("Hook response:", input.HookURL, "=>", res)

	return output, nil
}

func pullHook(hookURL string, output AnalyseOutput) (DoneHookResponse, error) {
	log.Println("Pulling hook:", hookURL)
	var response DoneHookResponse

	req, err := json.Marshal(map[string]string{
		"targetURL": output.TargetURL,
		"resultURL": output.ResultURL})
	if err != nil {
		return response, err
	}

	res, err := http.Post(hookURL, "application/json", bytes.NewBuffer(req))
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

func getResultURL(articleKey string) string {
	return ReadURL + "/" + articleKey
}
