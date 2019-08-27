package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// AnalyzeResult is the layout in which the final analysis are stored
type AnalyzeResult struct {
	Summary SummaryResult    `json:"summary"`
	Ranking RankPhraseResult `json:"ranking"`
}

func reduceArticle(articleKey string) (string, error) {
	log.Println("Reducing article:", articleKey)

	summaryResult, err := getSummaryResult(articleKey)
	if err != nil {
		return "", err
	}
	rankedResult, err := getRankedResult(articleKey)
	if err != nil {
		return "", err
	}

	result := AnalyzeResult{
		Summary: summaryResult,
		Ranking: rankedResult}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	doneFilePath := getArticleDonePath(articleKey)
	file, err := os.Create(doneFilePath)
	if err != nil {
		log.Println("Failed creating file:", doneFilePath, "=>", err.Error())
		return "", err
	}
	defer file.Close()

	file.Write(jsonResult)
	file.Sync()
	log.Println("Article reduced (done!):", doneFilePath)

	return doneFilePath, nil
}

func getSummaryResult(articleKey string) (SummaryResult, error) {
	var summary SummaryResult

	summarizedFilePath := getSummarizedArticlePath(articleKey)
	summaryBytes, err := ioutil.ReadFile(summarizedFilePath)
	if err != nil {
		return summary, err
	}

	err = json.Unmarshal(summaryBytes, &summary)
	if err != nil {
		return summary, err
	}

	return summary, nil
}

func getRankedResult(articleKey string) (RankPhraseResult, error) {
	var ranked RankPhraseResult

	rankedFilePath := getRankedArticlePath(articleKey)
	rankedBytes, err := ioutil.ReadFile(rankedFilePath)
	if err != nil {
		return ranked, err
	}

	err = json.Unmarshal(rankedBytes, &ranked)
	if err != nil {
		return ranked, err
	}

	return ranked, nil
}

func getArticleDonePath(articleKey string) string {
	return "data/" + articleKey + ".done"
}
