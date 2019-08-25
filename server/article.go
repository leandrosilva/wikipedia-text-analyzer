package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

func getArticleKey(targetURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(targetURL))
	hashed := hasher.Sum(nil)
	articleKey := fmt.Sprintf("%x", hashed)

	return articleKey
}

func getArticlePath(articleKey string, state string) string {
	return "data/" + articleKey + "." + state
}

func readArticleContent(articleKey string, state string) ([]byte, error) {
	filePath := getArticlePath(articleKey, state)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
