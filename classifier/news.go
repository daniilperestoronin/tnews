package classifier

import (
	"fmt"
	"strings"

	. "github.com/jbrukh/bayesian"
)

const (
	Good Class = "Good"
	Bad  Class = "Bad"
)

func NewsClassifier(tx string) {
	clsf := NewClassifierTfIdf(Good, Bad)
	goodStuff := []string{"tall", "rich", "handsome"}
	badStuff := []string{"poor", "smelly", "ugly"}
	clsf.Learn(goodStuff, Good)
	clsf.Learn(badStuff, Bad)
	clsf.ConvertTermsFreqToTfIdf()

	sc, _, _ := clsf.LogScores(strings.Split(tx, " "))

	fmt.Println(sc)
}
