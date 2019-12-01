package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/daniilperestoronin/nlp"
	"gonum.org/v1/gonum/mat"
)

const (
	enStopWrds      = "./corpus/en/stop_words"
	enNewsCorp      = "./corpus/en/news"
	enNotNewsCorp   = "./corpus/en/not_news"
	enSociety       = "./corpus/en/society"
	enEconomy       = "./corpus/en/economy"
	enTechnology    = "./corpus/en/technology"
	enSports        = "./corpus/en/sports"
	enEntertainment = "./corpus/en/entertainment"
	enScience       = "./corpus/en/science"
	enOther         = "./corpus/en/other"
	ruStopWrds      = "./corpus/ru/stop_words"
	ruNewsCorp      = "./corpus/ru/news"
	ruNotNewsCorp   = "./corpus/ru/not_news"
	ruSociety       = "./corpus/ru/society"
	ruEconomy       = "./corpus/ru/economy"
	ruTechnology    = "./corpus/ru/technology"
	ruSports        = "./corpus/ru/sports"
	ruEntertainment = "./corpus/ru/entertainment"
	ruScience       = "./corpus/ru/science"
	ruOther         = "./corpus/ru/other"
	tSVD            = 4
)

func main() {
	createNewsClassifier("./bin/en/", "news", enNewsCorp, enNotNewsCorp, enStopWrds, tSVD)
	createNewsGroupClassifier("./bin/en/", "newsGroup", enSociety, enEconomy, enTechnology, enSports, enEntertainment, enScience, enOther, enStopWrds, tSVD)

	createNewsClassifier("./bin/ru/", "news", ruNewsCorp, ruNotNewsCorp, ruStopWrds, tSVD)
	createNewsGroupClassifier("./bin/ru/", "newsGroup", ruSociety, ruEconomy, ruTechnology, ruSports, ruEntertainment, ruScience, ruOther, ruStopWrds, tSVD)
}

func readFileAsString(fileName string) string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func createNewsClassifier(folderPath string, fileName string, newsCorpPath string, notNewsCorpPath string, stopWordsPath string, tSVD int) {
	testCorpus := []string{
		readFileAsString(newsCorpPath),
		readFileAsString(notNewsCorpPath),
	}

	var stopWords = strings.Split(readFileAsString(stopWordsPath), "\n")

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()

	reducer := nlp.NewTruncatedSVD(tSVD)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	lsi, err := lsiPipeline.FitTransform(testCorpus...)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fo, err := os.Create(folderPath + fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(fo)

	dns := mat.DenseCopyOf(lsi)
	dns.MarshalBinaryTo(w)
	fmt.Println("create : " + fileName)
}

func createNewsGroupClassifier(folderPath string, fileName string, society string, economy string, technology string,
	sports string, entertainment string, science string, other string, stopWordsPath string, tSVD int) {

	testCorpus := []string{
		readFileAsString(society),
		readFileAsString(economy),
		readFileAsString(technology),
		readFileAsString(sports),
		readFileAsString(entertainment),
		readFileAsString(science),
		readFileAsString(other),
	}

	var stopWords = strings.Split(readFileAsString(stopWordsPath), "\n")

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()

	reducer := nlp.NewTruncatedSVD(tSVD)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	lsi, err := lsiPipeline.FitTransform(testCorpus...)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fo, err := os.Create(folderPath + fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(fo)

	dns := mat.DenseCopyOf(lsi)
	dns.MarshalBinaryTo(w)
	fmt.Println("create : " + fileName)
}
