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
	K                       int          `json:"k"`
	KMostImportantPhrases   []RankPhrase `json:"kMostImportantPhrases"`
	KMostImportantSentences []string     `json:"kMostImportantSentences"`
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

	result := RankPhraseResult{
		K: k,
		KMostImportantPhrases:   getKMostImportantPhrases(tr, k),
		KMostImportantSentences: getKMostImportantSentences(tr, k)}
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

func getKMostImportantPhrases(tr *textrank.TextRank, k int) []RankPhrase {
	var phrases []RankPhrase

	rankedPhrases := textrank.FindPhrases(tr)
	for i := range rankedPhrases {
		if i == k {
			break
		}
		phrase := rankedPhrases[i]
		phrases = append(phrases, RankPhrase{
			LeftText:   phrase.Left,
			RightText:  phrase.Right,
			Occurrence: phrase.Qty,
			Weight:     phrase.Weight})
	}

	return phrases
}

func getKMostImportantSentences(tr *textrank.TextRank, k int) []string {
	var sentences []string

	rankedSentences := textrank.FindSentencesByWordQtyWeight(tr, k)
	for i := range rankedSentences {
		if i == k {
			break
		}
		sentence := rankedSentences[i]
		sentences = append(sentences, sentence.Value)
	}

	return sentences
}

func getRankedArticlePath(articleKey string) string {
	return "data/" + articleKey + ".ranked"
}
