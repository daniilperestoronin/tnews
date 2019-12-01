package classifier

import (
	"fmt"

	goose "github.com/advancedlogic/GoOse"
	"github.com/daniilperestoronin/nlp"
	"github.com/daniilperestoronin/nlp/measures/pairwise"
	"gonum.org/v1/gonum/mat"
)

const (
	engStopWrds = "./alg/corpus/eng/stop_words"
	engNewsCorp = "./corpus/eng/news"
	tSVD        = 4
)

func NewsClassifier(tx string, lsi mat.Dense, stopWords []string) bool {

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()
	reducer := nlp.NewTruncatedSVD(tSVD)
	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	queryVector, err := lsiPipeline.FitTransform(tx)
	if err != nil {
		panic(err)
	}

	highestSimilarity := -1.0
	var matched int
	_, docs := lsi.Dims()
	for i := 0; i < docs; i++ {
		similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.ColView(i))
		if similarity > highestSimilarity {
			matched = i
			highestSimilarity = similarity
		}
	}

	return matched == 0
}

func NewsGroupClassifier(tx string, lsi mat.Dense, stopWords []string) string {

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()
	reducer := nlp.NewTruncatedSVD(tSVD)
	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	queryVector, err := lsiPipeline.FitTransform(tx)
	if err != nil {
		panic(err)
	}

	highestSimilarity := -1.0
	var matched int
	_, docs := lsi.Dims()
	for i := 0; i < docs; i++ {
		similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.ColView(i))
		if similarity > highestSimilarity {
			matched = i
			highestSimilarity = similarity
		}
	}

	switch matched {
	case 0:
		return "society"
	case 1:
		return "economy"
	case 2:
		return "technology"
	case 3:
		return "sports"
	case 4:
		return "entertainment"
	case 5:
		return "science"
	default:
		return "other"
	}
}

type Pair struct {
	Id         int
	Similarity float64
}

func NewsTreads(articles []*goose.Article, stopWords []string) {

	corpus := []string{}

	for _, article := range articles {
		corpus = append(corpus, article.CleanedText)
	}

	vectoriser := nlp.NewCountVectoriser(stopWords...)
	transformer := nlp.NewTfidfTransformer()

	// set k (the number of dimensions following truncation) to 4
	reducer := nlp.NewTruncatedSVD(100)

	lsiPipeline := nlp.NewPipeline(vectoriser, transformer, reducer)

	// Transform the corpus into an LSI fitting the model to the documents in the process
	lsi, err := lsiPipeline.FitTransform(corpus...)
	if err != nil {
		fmt.Printf("Failed to process documents because %v", err)
		return
	}

	aThread := map[int][]Pair{}

	for ai, article := range articles {
		aThread[ai] = []Pair{}
		queryVector, err := lsiPipeline.Transform(article.CleanedText)
		if err != nil {
			fmt.Printf("Failed to process documents because %v", err)
			return
		}
		_, docs := lsi.Dims()
		for i := 0; i < docs; i++ {
			similarity := pairwise.CosineSimilarity(queryVector.(mat.ColViewer).ColView(0), lsi.(mat.ColViewer).ColView(i))
			if similarity > 0.1 && i != ai {
				aThread[ai] = append(aThread[ai], Pair{Id: i, Similarity: similarity})
			}
		}
	}

	fmt.Println(aThread)
}
