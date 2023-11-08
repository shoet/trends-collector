package main

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_GoogleTrendsCrawler_FetchDocument(t *testing.T) {
	client := &http.Client{}
	sut := NewGoogleTrendsCrawler(client)

	url := "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja"
	doc, err := sut.FetchDocument(url)
	if err != nil {
		t.Fatalf("failed fetch document: %v", err)
	}

	html, err := doc.Html()
	if err != nil {
		t.Fatalf("failed get html: %v", err)
	}
	fmt.Println(html)

}
