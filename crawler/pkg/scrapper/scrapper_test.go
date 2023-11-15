package scrapper

import (
	"fmt"
	"testing"

	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector/util/timeutil"
)

func Test_GoogleTrendsDailyTrendsScrapper_ScrapePage(t *testing.T) {
	url := "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja"

	c := &timeutil.RealClocker{}
	browser, err := fetcher.BuildBrowser("/opt/homebrew/bin/chromium")
	if err != nil {
		t.Fatalf("failed build browser: %v", err)
	}

	doc := fetcher.FetchPage(browser, url)

	sut := NewGoogleTrendsDailyTrendsScrapper(c)
	pages, err := sut.ScrapePage("DailyTrend", doc)
	if err != nil {
		t.Fatalf("failed scrape page: %v", err)
	}

	fmt.Println(pages[0].TrendRank)

}

func Test_GoogleTrendsRealTimeTrendsScrapper_ScrapePage(t *testing.T) {
	url := "https://trends.google.co.jp/trends/trendingsearches/realtime?geo=JP&hl=ja&category=all"

	c := &timeutil.RealClocker{}
	browser, err := fetcher.BuildBrowser("/opt/homebrew/bin/chromium")
	if err != nil {
		t.Fatalf("failed build browser: %v", err)
	}

	doc := fetcher.FetchPage(browser, url)

	sut := NewGoogleTrendsRealTimeTrendsScrapper(c)
	pages, err := sut.ScrapePage("RealTimeTrend", doc)
	if err != nil {
		t.Fatalf("failed scrape page: %v", err)
	}

	fmt.Println(pages[0].TrendRank)
}
