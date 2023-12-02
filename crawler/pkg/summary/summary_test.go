package summary

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
)

func Test_MakeSummary(t *testing.T) {
	fetcher, err := fetcher.NewPageFetcher(&fetcher.PageFetcherInput{
		BrowserPath: "/opt/homebrew/bin/chromium",
	})
	if err != nil {
		t.Fatalf("failed to create fetcher: %v", err)
	}

	url := "https://www.yomiuri.co.jp/national/20231130-OYT1T50299/"
	// url := "https://news.yahoo.co.jp/articles/f8f54a0c5f6cb8db613b4e2be1c2f28ccf1342ee"
	// url := "https://news.yahoo.co.jp/articles/33e6c534669414b7028b02ba4384ad2ed9c5a622"

	summary, err := MakeSummary(fetcher, url)
	if err != nil {
		t.Fatalf("failed to make summary: %v", err)
	}

	fmt.Println(summary)

}

func Test_summaryTemplateBuilder(t *testing.T) {
	v := &OutputElement{
		Title:   "title",
		Content: "content",
	}

	summaryTemplate, err := summaryTemplateBuilder(v)
	if err != nil {
		t.Fatalf("failed to build template: %v", err)
	}

	if !strings.Contains(summaryTemplate, "###\ntitle\n###") {
		t.Fatalf("no match excepted title. want: %v, got: %v", "###\ntitle\n###", summaryTemplate)
	}

	if !strings.Contains(summaryTemplate, "###\ncontent\n###") {
		t.Fatalf("no match excepted content. want: %v, got: %v", "###\ncontent\n###", summaryTemplate)
	}
}
