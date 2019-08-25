package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

func getWikipediaArticleKey(targetURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(targetURL))
	hashed := hasher.Sum(nil)
	articleKey := fmt.Sprintf("%x", hashed)

	return articleKey
}

func downloadWikipediaArticle(targetURL string, articleKey string) (string, error) {
	log.Println("Downloading article:", targetURL)

	filePath := getWikipediaArticlePath(articleKey)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed creating file:", filePath, "=>", err.Error())
		return "", err
	}
	defer file.Close()

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		colly.MaxDepth(1),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*wikipedia.org",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		file.WriteString("[Source=" + targetURL + "]\n")
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		file.WriteString(e.Text)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// nothing here for now, maybe later on
	})

	c.Visit(targetURL)
	c.Wait()

	file.Sync()
	log.Println("Article downloaded:", targetURL, "=>", filePath)

	return filePath, nil
}

func readWikipediaRawArticle(articleKey string) ([]byte, error) {
	filePath := getWikipediaArticlePath(articleKey)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func getWikipediaArticlePath(articleKey string) string {
	return "data/" + articleKey + ".txt"
}
