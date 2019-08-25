package main

import "log"

func clearDocument(filePath string) (string, error) {
	log.Println("Cleaning article:", filePath)
	log.Println("Article cleaned:", filePath)

	return "blah", nil
}
