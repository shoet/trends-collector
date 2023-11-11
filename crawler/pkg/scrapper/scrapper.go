package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/shoet/trends-collector/entities"
)

type Scrapper interface {
	ScrapePage(string, *rod.Page) ([]*entities.Page, error)
}

type Scrappers []struct {
	Category string
	Url      string
	Scrapper Scrapper
	Pages    []*entities.Page
}

type GoogleTrendsDailyTrendsScrapper struct{}

func (g *GoogleTrendsDailyTrendsScrapper) ScrapePage(
	category string,
	doc *rod.Page,
) ([]*entities.Page, error) {
	topElement := doc.MustElement(".feed-list-wrapper")
	trends := topElement.MustElements(".md-list-block")

	pages := make([]*entities.Page, 0, len(trends))
	for _, element := range trends {
		title := element.MustElement(".title").MustText()
		count := element.MustElement(".search-count-title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		pages = append(pages, &entities.Page{
			Category: category,
			Url:      *url,
			Title:    title,
			Trend:    count,
			// TODO: partition yyyymmdd
		})
	}
	return pages, nil
}

type GoogleTrendsRealTimeTrendsScrapper struct{}

func (g *GoogleTrendsRealTimeTrendsScrapper) ScrapePage(
	category string,
	doc *rod.Page,
) ([]*entities.Page, error) {
	topElement := doc.MustElement(".trending-feed-content")
	trends := topElement.MustElements(".feed-item-header")

	pages := make([]*entities.Page, 0, len(trends))
	for _, element := range trends {
		title := element.MustElement(".title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		pages = append(pages, &entities.Page{
			Category: category,
			Url:      *url,
			Title:    title,
			// TODO: partition yyyymmddHHMMSS
		})
	}
	return pages, nil
}
