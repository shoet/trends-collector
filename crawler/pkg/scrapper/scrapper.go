package scrapper

import (
	"github.com/go-rod/rod"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util/timeutil"
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

type GoogleTrendsDailyTrendsScrapper struct {
	clocker interfaces.Clocker
}

func NewGoogleTrendsDailyTrendsScrapper(
	clocker interfaces.Clocker,
) *GoogleTrendsDailyTrendsScrapper {
	return &GoogleTrendsDailyTrendsScrapper{clocker: clocker}
}

func (g *GoogleTrendsDailyTrendsScrapper) ScrapePage(
	category string,
	doc *rod.Page,
) ([]*entities.Page, error) {
	topElement := doc.MustElement(".feed-list-wrapper")
	trends := topElement.MustElements(".md-list-block")
	pages := make([]*entities.Page, 0, len(trends))
	ymd := timeutil.NowFormatYYYYMMDD(g.clocker)
	for _, element := range trends {
		title := element.MustElement(".title").MustText()
		count := element.MustElement(".search-count-title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		pages = append(pages, &entities.Page{
			Partition: ymd,
			Category:  category,
			Url:       *url,
			Title:     title,
			Trend:     count,
		})
	}
	return pages, nil
}

type GoogleTrendsRealTimeTrendsScrapper struct {
	clocker interfaces.Clocker
}

func NewGoogleTrendsRealTimeTrendsScrapper(
	clocker interfaces.Clocker,
) *GoogleTrendsRealTimeTrendsScrapper {
	return &GoogleTrendsRealTimeTrendsScrapper{clocker: clocker}
}

func (g *GoogleTrendsRealTimeTrendsScrapper) ScrapePage(
	category string,
	doc *rod.Page,
) ([]*entities.Page, error) {
	topElement := doc.MustElement(".trending-feed-content")
	trends := topElement.MustElements(".feed-item-header")
	pages := make([]*entities.Page, 0, len(trends))
	ymdhms := timeutil.NowFormatYYYYMMDDHHMMSS(g.clocker)
	for _, element := range trends {
		title := element.MustElement(".title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		pages = append(pages, &entities.Page{
			Partition: ymdhms,
			Category:  category,
			Url:       *url,
			Title:     title,
		})
	}
	return pages, nil
}
