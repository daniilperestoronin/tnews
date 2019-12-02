package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	goose "github.com/advancedlogic/GoOse"
	"github.com/daniilperestoronin/tnews/classifier"
	"github.com/daniilperestoronin/tnews/lang"
	"github.com/daniilperestoronin/tnews/parse"
	"github.com/urfave/cli"
	"gonum.org/v1/gonum/mat"
)

const (
	enStopWrds       = "./alg/corpus/en/stop_words"
	enNewsDense      = "./alg/bin/en/news"
	enNewsGroupDense = "./alg/bin/en/newsGroup"
	ruStopWrds       = "./alg/corpus/ru/stop_words"
	ruNewsDense      = "./alg/bin/ru/news"
	ruNewsGroupDense = "./alg/bin/ru/newsGroup"
)

func cliApp() *cli.App {

	nlpModels := loadNlpModels()

	app := cli.NewApp()
	app.Name = "tnews"
	app.Usage = "news clustering"
	app.Version = "0.0.1"
	app.UsageText = "tnews [comand] source_dir"
	app.Commands = []cli.Command{
		{
			Name:  "languages",
			Usage: "Isolate articles in english and Russian",
			Action: func(c *cli.Context) error {
				printResultJSON(checkLanguages(getSrcDir(c)))
				return nil
			},
		},
		{
			Name:  "news",
			Usage: "Isolate news articles",
			Action: func(c *cli.Context) error {
				printResultJSON(checkNews(getSrcDir(c), nlpModels))
				return nil
			},
		},
		{
			Name:  "categories",
			Usage: "Group news articles by category",
			Action: func(c *cli.Context) error {
				printResultJSON(checkNewsGroup(getSrcDir(c), nlpModels))
				return nil
			},
		},
		{
			Name:  "threads",
			Usage: "Group similar news into threads",
			Action: func(c *cli.Context) error {
				printResultJSON(checkNewsTreads(getSrcDir(c), nlpModels))
				return nil
			},
		},
		{
			Name:  "top",
			Usage: "Sort threads by their relative importance",
			Action: func(c *cli.Context) error {
				printResultJSON(checkNewsTreadsByCategory(getSrcDir(c), nlpModels))
				return nil
			},
		},
	}
	return app
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
		panic(err)
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

func checkNews(filesPath string, nlpModels map[string]nlpModel) newsArr {
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
				aLang := lang.DetectLanguage(article.Title + article.CleanedText)
				if aLang == "en" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["en"].News, nlpModels["en"].StopWords) {
						articles = append(articles, info.Name())
					}
				} else if aLang == "ru" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["ru"].News, nlpModels["ru"].StopWords) {
						articles = append(articles, info.Name())
					}
				}
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	return newsArr{Articles: articles}
}

type newsGroup struct {
	Group    string   `json:"category"`
	Articles []string `json:"articles"`
}

func checkNewsGroup(filesPath string, nlpModels map[string]nlpModel) []newsGroup {
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
				aLang := lang.DetectLanguage(article.Title + article.CleanedText)
				if aLang == "en" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["en"].News, nlpModels["en"].StopWords) {
						group := classifier.NewsGroupClassifier(article.Title+article.CleanedText, nlpModels["en"].NewsGroup, nlpModels["en"].StopWords)
						nGroup[group] = append(nGroup[group], info.Name())
					}
				} else if aLang == "ru" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["ru"].News, nlpModels["ru"].StopWords) {
						group := classifier.NewsGroupClassifier(article.Title+article.CleanedText, nlpModels["ru"].NewsGroup, nlpModels["ru"].StopWords)
						nGroup[group] = append(nGroup[group], info.Name())
					}
				}
			}
			return nil
		})

	if err != nil {
		panic(err)
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

func checkNewsTreads(filesPath string, nlpModels map[string]nlpModel) []newsThread {

	enArticles := []*goose.Article{}
	ruArticles := []*goose.Article{}

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
				aLang := lang.DetectLanguage(article.Title + article.CleanedText)
				if aLang == "en" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["en"].News, nlpModels["en"].StopWords) {
						article.FinalURL = info.Name()
						enArticles = append(enArticles, article)
					}
				} else if aLang == "ru" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["ru"].News, nlpModels["ru"].StopWords) {
						article.FinalURL = info.Name()
						ruArticles = append(ruArticles, article)
					}
				}
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	newsThreads := []newsThread{}

	newsThreads = append(newsThreads, convertToNewsThread(classifier.NewsTreads(enArticles, nlpModels["en"].StopWords))...)
	newsThreads = append(newsThreads, convertToNewsThread(classifier.NewsTreads(ruArticles, nlpModels["ru"].StopWords))...)

	return newsThreads
}

func convertToNewsThread(nThr map[string][]string) []newsThread {
	newsThreads := []newsThread{}
	for k, v := range nThr {
		newsThreads = append(newsThreads, newsThread{Title: k, Articles: v})
	}
	return newsThreads
}

type newsGroupTreads struct {
	Group      string       `json:"category"`
	NewsThread []newsThread `json:"threads"`
}

func checkNewsTreadsByCategory(filesPath string, nlpModels map[string]nlpModel) []newsGroupTreads {

	enNGroup := map[string][]*goose.Article{}
	ruNGroup := map[string][]*goose.Article{}

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
				aLang := lang.DetectLanguage(article.Title + article.CleanedText)
				if aLang == "en" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["en"].News, nlpModels["en"].StopWords) {
						article.FinalURL = info.Name()
						group := classifier.NewsGroupClassifier(article.Title+article.CleanedText, nlpModels["en"].NewsGroup, nlpModels["en"].StopWords)
						enNGroup[group] = append(enNGroup[group], article)
					}
				} else if aLang == "ru" {
					if classifier.NewsClassifier(article.Title+article.CleanedText, nlpModels["ru"].News, nlpModels["ru"].StopWords) {
						article.FinalURL = info.Name()
						group := classifier.NewsGroupClassifier(article.Title+article.CleanedText, nlpModels["ru"].NewsGroup, nlpModels["ru"].StopWords)
						ruNGroup[group] = append(ruNGroup[group], article)
					}
				}
			}
			return nil
		})

	if err != nil {
		panic(err)
	}

	l := []newsGroupTreads{}

	for k, v := range enNGroup {

		nThreads := []newsThread{}

		nThreads = append(nThreads, convertToNewsThread(classifier.NewsTreads(v, nlpModels["en"].StopWords))...)

		if ruNGroup[k] != nil {
			nThreads = append(nThreads, convertToNewsThread(classifier.NewsTreads(ruNGroup[k], nlpModels["ru"].StopWords))...)
			delete(ruNGroup, k)
		}

		l = append(l,
			newsGroupTreads{
				Group:      k,
				NewsThread: nThreads,
			})
	}

	for k, v := range ruNGroup {

		l = append(l,
			newsGroupTreads{
				Group:      k,
				NewsThread: convertToNewsThread(classifier.NewsTreads(v, nlpModels["ru"].StopWords)),
			})
	}

	return l
}

type nlpModel struct {
	StopWords []string
	News      mat.Dense
	NewsGroup mat.Dense
}

func loadNlpModels() map[string]nlpModel {
	return map[string]nlpModel{
		"en": nlpModel{
			StopWords: strings.Split(readFileAsString(enStopWrds), "\n"),
			News:      getDenseFromBin(enNewsDense),
			NewsGroup: getDenseFromBin(enNewsGroupDense),
		},
		"ru": nlpModel{
			StopWords: strings.Split(readFileAsString(ruStopWrds), "\n"),
			News:      getDenseFromBin(ruNewsDense),
			NewsGroup: getDenseFromBin(ruNewsGroupDense),
		},
	}
}

func getDenseFromBin(fileName string) mat.Dense {
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	r := bufio.NewReader(fi)
	lsi := mat.Dense{}
	lsi.UnmarshalBinaryFrom(r)
	return lsi
}

func readFileAsString(fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func getSrcDir(c *cli.Context) string {
	srcDir := c.Args().First()
	if srcDir == "" {
		panic("provide source_dir")
	}
	return srcDir
}

func printResultJSON(n interface{}) {
	prettyJSON, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}
