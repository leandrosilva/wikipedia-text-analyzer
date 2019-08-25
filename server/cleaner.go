package main

import "log"

func cleanArticle(articleKey string, filePath string) (string, error) {
	log.Println("Cleaning article:", filePath)

	cleanedFilePath := getCleanedArticlePath(articleKey)
	log.Println("Article cleaned:", cleanedFilePath)

	return cleanedFilePath, nil
}

func getCleanedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".cleaned"
}
