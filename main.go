package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/daniilperestoronin/tnews/lang"
	"github.com/daniilperestoronin/tnews/parse"
)

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
				article := parseArticleFromHTMLFile(string(b))
				lang := detectLanguage(article.Title + article.CleanedText)
				fmt.Println(path + " - " + lang)
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}

	app := cliApp()

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
