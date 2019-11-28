package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/daniilperestoronin/tnews/classifier"
	"github.com/daniilperestoronin/tnews/lang"
	"github.com/daniilperestoronin/tnews/parse"
	"github.com/urfave/cli"
)

func cliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "tnews"
	app.Usage = "news clustering"
	app.Version = "0.0.1"
	app.UsageText = "tnews [comand] source_dir"
	app.Commands = []cli.Command{
		{
			Name:  "languages",
			Usage: "Isolate articles in English and Russian",
			Action: func(c *cli.Context) error {
				printResultJson(checkLanguages(getSrcDir(c)))
				return nil
			},
		},
		{
			Name:  "news",
			Usage: "Isolate news articles",
			Action: func(c *cli.Context) error {
				printResultJson(checkNews(getSrcDir(c)))
				return nil
			},
		},
		{
			Name:  "categories",
			Usage: "Group news articles by category",
			Action: func(c *cli.Context) error {
				printResultJson(checkNewsGroup(getSrcDir(c)))
				return nil
			},
		},
		{
			Name:  "threads",
			Usage: "Group similar news into threads",
			Action: func(c *cli.Context) error {
				printResultJson(checkNewsTreads(getSrcDir(c)))
				return nil
			},
		},
		{
			Name:  "top",
			Usage: "Sort threads by their relative importance",
			Action: func(c *cli.Context) error {
				printResultJson(checkNewsTreadsByCategory(getSrcDir(c)))
				return nil
			},
		},
	}
	return app
}

func getSrcDir(c *cli.Context) string {
	srcDir := c.Args().First()
	if srcDir == "" {
		panic("provide source_dir")
	}
	return srcDir
}

func printResultJson(n interface{}) {
	prettyJSON, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}

type checkLanguagesStr struct {
	LangCode string   `json:"lang_code"`
	Articles []string `json:"articles"`
}

func checkLanguages(filesPath string) []checkLanguagesStr {
	aLang := make(map[string][]string)

	err := filepath.Walk(filesPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				article := parse.ParseArticleFromHTMLFile(string(b))
				lang := lang.DetectLanguage(article.Title + article.CleanedText)
				aLang[lang] = append(aLang[lang], info.Name())
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}

	l := []checkLanguagesStr{}

	l = append(l, checkLanguagesStr{LangCode: "en", Articles: aLang["en"]})
	delete(aLang, "en")
	l = append(l, checkLanguagesStr{LangCode: "ru", Articles: aLang["ru"]})
	delete(aLang, "ru")

	for k, v := range aLang {
		l = append(l, checkLanguagesStr{LangCode: k, Articles: v})
	}

	return l
}

type newsArr struct {
	Articles []string `json:"articles"`
}

func checkNews(filesPath string) newsArr {
	articles := []string{}
	err := filepath.Walk(filesPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				article := parse.ParseArticleFromHTMLFile(string(b))
				if lang.DetectLanguage(article.Title+article.CleanedText) == "en" {
					if classifier.NewsClassifier(article.Title + article.CleanedText) {
						articles = append(articles, info.Name())
					}
				}
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}

	return newsArr{Articles: articles}
}

type newsGroup struct {
	Group    string   `json:"category"`
	Articles []string `json:"articles"`
}

func checkNewsGroup(filesPath string) []newsGroup {
	nGroup := make(map[string][]string)

	err := filepath.Walk(filesPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				b, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}
				article := parse.ParseArticleFromHTMLFile(string(b))
				if lang.DetectLanguage(article.Title+article.CleanedText) == "en" {
					if classifier.NewsClassifier(article.Title + article.CleanedText) {
						group := classifier.NewsGroupClassifier(article.Title + article.CleanedText)
						nGroup[group] = append(nGroup[group], info.Name())
					}
				}
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}

	l := []newsGroup{}

	for k, v := range nGroup {
		l = append(l, newsGroup{Group: k, Articles: v})
	}

	return l
}

type newsThread struct {
	Title    string   `json:"title"`
	Articles []string `json:"articles"`
}

func checkNewsTreads(filesPath string) []newsThread {
	return nil
}

func checkNewsTreadsByCategory(filesPath string) []newsThread {
	return nil
}
