package classifier

import (
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
