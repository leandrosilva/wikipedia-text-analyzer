package main

import "log"

func summarizeArticle(articleKey string, rawFilePath string) (string, error) {
	log.Println("Summarizing article:", rawFilePath)

	summarizedFilePath := getSummarizedArticlePath(articleKey)
	log.Println("Article summarized:", summarizedFilePath)

	return summarizedFilePath, nil
}

func getSummarizedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".summarized"
}
