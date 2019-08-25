package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	textrank "github.com/DavidBelicza/TextRank"
)

// RankPhrase is a pair of words occurence
type RankPhrase struct {
	LeftText   string  `json:"leftText"`
	RightText  string  `json:"rightText"`
	Occurrence int     `json:"occurrence"`
	Weight     float32 `json:"weight"`
}

// RankPhraseResult is the final product of this phrase analysis
type RankPhraseResult struct {
	K                     int          `json:"k"`
	KMostImportantPhrases []RankPhrase `json:"kMostImportantPhrases"`
}

func rankArticlePhrases(articleKey string, rawFilePath string, k int) (string, error) {
	log.Println("Ranking article phrases:", rawFilePath)

	contentBytes, err := ioutil.ReadFile(rawFilePath)
	if err != nil {
		return "", err
	}

	rawText := string(contentBytes)
	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()

	tr.Populate(rawText, language, rule)
	tr.Ranking(algorithmDef)

	rankedPhrases := textrank.FindPhrases(tr)
	result := RankPhraseResult{K: k}

	for i := range rankedPhrases {
		if i == k {
			break
		}
		result.KMostImportantPhrases = append(result.KMostImportantPhrases, RankPhrase{
			LeftText:   rankedPhrases[i].Left,
			RightText:  rankedPhrases[i].Right,
			Occurrence: rankedPhrases[i].Qty,
			Weight:     rankedPhrases[i].Weight})
	}

	res, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	rankedFilePath := getRankedArticlePath(articleKey)
	file, err := os.Create(rankedFilePath)
	if err != nil {
		log.Println("Failed creating file:", rankedFilePath, "=>", err.Error())
		return "", err
	}
	defer file.Close()

	file.Write(res)
	file.Sync()
	log.Println("Article phrases ranked:", rankedFilePath)

	return rankedFilePath, nil
}

func getRankedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".ranked"
}
