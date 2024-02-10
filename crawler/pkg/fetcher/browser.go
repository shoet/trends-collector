package fetcher

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type RodPageFetcher struct {
	browser *rod.Browser
}

type PageFetcherInput struct {
	BrowserPath string
}

type FetchPageResult struct {
	RodPage *rod.Page
}

func NewRodPageFetcher(input *PageFetcherInput) (*RodPageFetcher, error) {
	if input.BrowserPath == "" {
		return nil, fmt.Errorf("Browser path is empty")
	}
	browser, err := BuildBrowser(input.BrowserPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to build browser: %w", err)
	}
	return &RodPageFetcher{browser: browser}, nil
}

func (f *RodPageFetcher) FetchPage(url string) (*FetchPageResult, error) {
	page := f.browser.MustPage(url)
	page.MustWaitLoad()

	result := &FetchPageResult{
		RodPage: page,
	}
	return result, nil
}

func BuildBrowser(browserPath string) (*rod.Browser, error) {
	u := launcher.New().Bin(browserPath).NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	fmt.Printf("Start browser: %s\n", u)
	return browser, nil
}
