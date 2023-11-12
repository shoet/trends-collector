package fetcher

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func BuildBrowser(browserPath string) (*rod.Browser, error) {
	u := launcher.New().Bin(browserPath).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	fmt.Printf("Start browser: %s\n", u)
	return browser, nil
}

func FetchPage(browser *rod.Browser, url string) *rod.Page {
	page := browser.MustPage(url)
	page.MustWaitLoad()
	return page
}
