package scrapper

import (
	"testing"
)

func Test_GoogleTrendsDailyTrendsScrapper_ScrapePage(t *testing.T) {
	t.Skip("refactoring")

	// client := &http.Client{}
	//
	// url := "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja"
	// c := webcrawler.NewWebCrawler()
	//
	// doc, err := c.GetPageHTML(url)
	// if err != nil {
	// 	t.Fatalf("failed create document: %v", err)
	// }
	//
	// h, err := doc.Html()
	// os.Stdout.WriteString(h)
	//
	// pages, err := sut.ScrapePage(doc)
	// if err != nil {
	// 	t.Fatalf("failed scrape page: %v", err)
	// }
	// fmt.Println(pages)

}

func Test_NewHHKBStudioNotifyScrapper(t *testing.T) {
	t.Skip("refactoring")
	// browserPath := "/opt/homebrew/bin/chromium"
	// browser, err := fetcher.BuildBrowser(browserPath)
	// if err != nil {
	// 	t.Fatalf("failed build browser: %v", err)
	// }
	// page := fetcher.FetchPage(
	// 	browser, "https://www.pfu.ricoh.com/direct/hhkb/hhkb-studio/detail_pd-id120b.html")
	//
	// s := NewHHKBStudioNotifyScrapper()
	// p, err := s.ScrapePage("HHKB", page)
	// if err != nil {
	// 	t.Fatalf("failed scrape page: %v", err)
	// }
	//
	// _ = p
}
