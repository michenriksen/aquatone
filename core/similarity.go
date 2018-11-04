package core

import (
	"io"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/html"
)

func GetPageStructure(body io.Reader) ([]string, error) {
	var structure []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return structure, nil
		case html.StartTagToken:
			tn, _ := z.TagName()
			structure = append(structure, string(tn))
		}
	}
}

func GetSimilarity(a, b []string) float64 {
	matcher := difflib.NewMatcher(a, b)
	return matcher.Ratio()
}
