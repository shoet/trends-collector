package fetcher

import (
	"github.com/go-rod/rod"
	"github.com/playwright-community/playwright-go"
)

type PageFetcherInput struct {
	BrowserPath string
}

type FetchPageResult struct {
	RodPage        *rod.Page
	PlaywrightPage playwright.Page
}
