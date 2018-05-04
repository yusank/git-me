package utils

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetDoc return Document object of the HTML string
func GetDoc(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return doc, nil
}

// Title get title
func Title(doc *goquery.Document) string {
	var title string
	title = strings.Replace(
		strings.TrimSpace(doc.Find("h1").First().Text()), "\n", "", -1,
	)
	if title == "" {
		// Bilibili: Some movie page got no h1 tag
		title, _ = doc.Find("meta[property=\"og:title\"]").Attr("content")
	}
	return FileName(title)
}
