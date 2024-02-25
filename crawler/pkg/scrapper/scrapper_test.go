package scrapper

import (
	"testing"

	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector/util/timeutil"
)

func Test_GoogleTrendsDailyTrendsScrapper_ScrapePage(t *testing.T) {
	url := "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja"

	c := &timeutil.RealClocker{}
	fetcher, closer, err := fetcher.NewRodPageFetcher(&fetcher.PageFetcherInput{
		BrowserPath: "/opt/homebrew/bin/chromium",
	})
	if err != nil {
		t.Fatalf("failed build fetcher: %v", err)
	}
	defer closer()

	result, err := fetcher.FetchPage(url)
	if err != nil {
		t.Fatalf("failed fetch page: %v", err)
	}

	sut := NewGoogleTrendsDailyTrendsScrapper(c)
	pages, err := sut.ScrapePage("DailyTrend", &ScrapperInput{
		RodPage: result.RodPage,
	})
	if err != nil {
		t.Fatalf("failed scrape page: %v", err)
	}

	if len(pages) == 0 {
		t.Fatalf("no pages")
	}

}

func Test_GoogleTrendsRealTimeTrendsScrapper_ScrapePage(t *testing.T) {
	url := "https://trends.google.co.jp/trends/trendingsearches/realtime?geo=JP&hl=ja&category=all"

	c := &timeutil.RealClocker{}
	fetcher, closer, err := fetcher.NewRodPageFetcher(&fetcher.PageFetcherInput{
		BrowserPath: "/opt/homebrew/bin/chromium",
	})
	defer closer()
	if err != nil {
		t.Fatalf("failed build fetcher: %v", err)
	}

	result, err := fetcher.FetchPage(url)
	if err != nil {
		t.Fatalf("failed fetch page: %v", err)
	}

	sut := NewGoogleTrendsRealTimeTrendsScrapper(c)
	pages, err := sut.ScrapePage("RealTimeTrend", &ScrapperInput{
		RodPage: result.RodPage,
	})
	if err != nil {
		t.Fatalf("failed scrape page: %v", err)
	}

	if len(pages) == 0 {
		t.Fatalf("no pages")
	}
}
