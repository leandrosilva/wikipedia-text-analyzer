package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/JesusIslam/tldr"
)

// SummaryResult is the layout in which text summarized are stored
type SummaryResult struct {
	K         int      `json:"k"`
	Sentences []string `json:"sentences"`
}

func summarizeArticle(articleKey string, rawFilePath string, k int) (string, error) {
	log.Println("Summarizing article:", rawFilePath)

	contentBytes, err := ioutil.ReadFile(rawFilePath)
	if err != nil {
		return "", err
	}

	rawText := string(contentBytes)
	bag := tldr.New()
	summary, _ := bag.Summarize(rawText, k)

	summarizedFilePath := getSummarizedArticlePath(articleKey)
	file, err := os.Create(summarizedFilePath)
	if err != nil {
		log.Println("Failed creating file:", summarizedFilePath, "=>", err.Error())
		return "", err
	}
	defer file.Close()

	result := SummaryResult{K: k}
	for i := range summary {
		result.Sentences = append(result.Sentences, summary[i])
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	file.Write(jsonResult)
	file.Sync()
	log.Println("Article summarized:", summarizedFilePath)

	return summarizedFilePath, nil
}

func getSummarizedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".summarized"
}
