package parse

import goose "github.com/advancedlogic/GoOse"

func ParseArticleFromHTMLFile(htmlText string) *goose.Article {
	g := goose.New()
	a, err := g.ExtractFromRawHTML(htmlText, "")
	if err != nil {
		panic("parseArticleFromHTMLFile error")
	}
	return a
}
