package main

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/microcosm-cc/bluemonday"
	"github.com/shoet/trends-collector/entities"
)

type GoogleTrendsDailyTrendsScrapper struct{}

func NewGoogleTrendsDailyTrendsScrapper() (*GoogleTrendsDailyTrendsScrapper, error) {
	return &GoogleTrendsDailyTrendsScrapper{}, nil
}

func (g *GoogleTrendsDailyTrendsScrapper) ScrapePage(
	doc *rod.Page,
) ([]*entities.Page, error) {

	topElement := doc.MustElement(".feed-list-wrapper")

	trends := topElement.MustElements(".md-list-block")

	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class", "id").OnElements("div")

	for _, element := range trends {
		title := element.MustElement(".title").MustText()
		count := element.MustElement(".search-count-title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		fmt.Println(title, count, *url)
	}
	return nil, nil
}

func (g *GoogleTrendsDailyTrendsScrapper) ScrapePage2(
	url string,
) ([]*entities.Page, error) {

	return nil, nil
}
