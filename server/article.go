package main

import "io/ioutil"

func readArticleContent(articleKey string, state string) ([]byte, error) {
	filePath := getArticlePath(articleKey, state)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func getArticlePath(articleKey string, state string) string {
	return "data/" + articleKey + "." + state
}
