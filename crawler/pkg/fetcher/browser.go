package fetcher

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type PageFetcher struct {
	browser *rod.Browser
}

type PageFetcherInput struct {
	BrowserPath string
}

func NewPageFetcher(input *PageFetcherInput) (*PageFetcher, error) {
	if input.BrowserPath == "" {
		return nil, fmt.Errorf("Browser path is empty")
	}
	browser, err := BuildBrowser(input.BrowserPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to build browser: %w", err)
	}
	return &PageFetcher{browser: browser}, nil
}

func (f *PageFetcher) FetchPage(url string) *rod.Page {
	page := f.browser.MustPage(url)
	page.MustWaitLoad()
	return page
}

func BuildBrowser(browserPath string) (*rod.Browser, error) {
	u := launcher.New().Bin(browserPath).NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	fmt.Printf("Start browser: %s\n", u)
	return browser, nil
}
