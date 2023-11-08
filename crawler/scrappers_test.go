package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func Test_GoogleTrendsDailyTrendsScrapper_ScrapePage(t *testing.T) {
	sut, err := NewGoogleTrendsDailyTrendsScrapper()
	if err != nil {
		t.Fatalf("failed create scrapper: %v", err)
	}

	client := &http.Client{}

	url := "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja"
	c := NewGoogleTrendsCrawler(client)

	doc, err := c.GetPageHTML(url)
	if err != nil {
		t.Fatalf("failed create document: %v", err)
	}

	h, err := doc.Html()
	os.Stdout.WriteString(h)

	pages, err := sut.ScrapePage(doc)
	if err != nil {
		t.Fatalf("failed scrape page: %v", err)
	}
	fmt.Println(pages)

}
