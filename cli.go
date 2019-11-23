package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
				srcDir := getSrcDir(c)
				l := checkLanguages(srcDir)
				fmt.Println(l)
				return nil
			},
		},
		{
			Name:  "news",
			Usage: "Isolate news articles",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				return nil
			},
		},
		{
			Name:  "categories",
			Usage: "Group news articles by category",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				return nil
			},
		},
		{
			Name:  "threads",
			Usage: "Group similar news into threads",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				return nil
			},
		},
		{
			Name:  "top",
			Usage: "Sort threads by their relative importance",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				return nil
			},
		},
	}

	return app
}

func getSrcDir(c *cli.Context) string {
	srcDir := c.Args().First()
	if srcDir == "" {
		fmt.Println("provide source_dir")
	}
	return srcDir
}

type checkLanguagesStr struct {
	langCode string   `json:"lang_code"`
	articles []string `json:"articles"`
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

	for k, v := range aLang {
		l = append(l, checkLanguagesStr{langCode: k, articles: v})
	}

	return l
}
