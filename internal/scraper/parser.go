package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// StripBoilerplate removes elements that add token cost without pricing signal:
// nav, footer, script, style, header, noscript, and aria-hidden elements.
// Returns the inner HTML of <body>.
func StripBoilerplate(rawHTML string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}
	doc.Find("nav").Remove()
	doc.Find("footer").Remove()
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	doc.Find("header").Remove()
	doc.Find("noscript").Remove()
	doc.Find("[aria-hidden='true']").Remove()

	stripped, err := doc.Find("body").Html()
	if err != nil {
		return "", err
	}
	return stripped, nil
}
