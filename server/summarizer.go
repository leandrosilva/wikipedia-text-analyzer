package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/JesusIslam/tldr"
)

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

	for i := range summary {
		file.WriteString(summary[i] + "\n\n")
	}

	file.Sync()
	log.Println("Article summarized:", summarizedFilePath)

	return summarizedFilePath, nil
}

func getSummarizedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".summarized"
}
