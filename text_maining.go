package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/abadojack/whatlanggo"
	goose "github.com/advancedlogic/GoOse"
)

func detectLanguage(text string) string {
	info := whatlanggo.Detect(text)
	return info.Lang.String()
}

func parceArticleFromHTMLFile(htmlText string) *goose.Article {
	g := goose.New()
	a, err := g.ExtractFromRawHTML(htmlText, "")
	if err != nil {
		panic("parceArticleFromHTMLFile error")
	}
	return a
}

func main() {
	err := filepath.Walk("../",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				article := parceArticleFromHTMLFile(string(b))
				lang := detectLanguage(article.Title + article.CleanedText)
				fmt.Println(path + " - " + lang)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
