package classifier

import (
	"fmt"
	"strings"

	"github.com/jbrukh/bayesian"
)

const (
	News    bayesian.Class = "News"
	NotNews bayesian.Class = "NotNews"
)

func NewsClassifier(tx string) {
	clsf := bayesian.NewClassifierTfIdf(News, NotNews)

	newsWords := []string{"tall", "rich", "handsome"}
	notNewsWords := []string{"poor", "smelly", "ugly"}

	clsf.Learn(newsWords, News)
	clsf.Learn(notNewsWords, NotNews)

	clsf.ConvertTermsFreqToTfIdf()

	_, i, _ := clsf.LogScores(strings.Split(tx, " "))

	fmt.Println(i)
}

const (
	//Society (includes Politics, Elections, Legislation, Incidents, Crime)
	Society bayesian.Class = "Society"
	//Economy (includes Markets, Finance, Business)
	Economy bayesian.Class = "Economy"
	//Technology (includes Gadgets, Auto, Apps, Internet services)
	Technology bayesian.Class = "Technology"
	//Sports (includes E-Sports)
	Sports bayesian.Class = "Sports"
	//Entertainment (includes Movies, Music, Games, Books, Arts)
	Entertainment bayesian.Class = "Entertainment"
	//Science (includes Health, Biology, Physics, Genetics)
	Science bayesian.Class = "Science"
	//Other (news articles that don't fall into any of the above categories)
	Other bayesian.Class = "Other"
)

func NewsGroupClassifier(tx string) {
	clsf := bayesian.NewClassifierTfIdf(Society, Economy, Technology, Sports, Entertainment, Science, Other)

	societyWords := []string{"tall", "rich", "handsome"}
	economyWords := []string{"tall", "rich", "handsome"}
	technologyWords := []string{"tall", "rich", "handsome"}
	sportsWords := []string{"tall", "rich", "handsome"}
	entertainmentWords := []string{"tall", "rich", "handsome"}
	scienceWords := []string{"tall", "rich", "handsome"}
	otherWords := []string{"tall", "rich", "handsome"}

	clsf.Learn(societyWords, Society)
	clsf.Learn(economyWords, Economy)
	clsf.Learn(technologyWords, Technology)
	clsf.Learn(sportsWords, Sports)
	clsf.Learn(entertainmentWords, Entertainment)
	clsf.Learn(scienceWords, Science)
	clsf.Learn(otherWords, Other)

	clsf.ConvertTermsFreqToTfIdf()

	_, i, _ := clsf.LogScores(strings.Split(tx, " "))

	fmt.Println(i)
}
