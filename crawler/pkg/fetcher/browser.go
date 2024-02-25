package fetcher

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type RodPageFetcher struct {
	browser *rod.Browser
}

func NewRodPageFetcher(input *PageFetcherInput) (*RodPageFetcher, func() error, error) {
	if input.BrowserPath == "" {
		return nil, nil, fmt.Errorf("BrowserPath is required")
	}
	browser, cleanup, err := BuildBrowser(input.BrowserPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to build browser: %w", err)
	}
	return &RodPageFetcher{browser: browser}, cleanup, nil
}

func (f *RodPageFetcher) FetchPage(url string) (*FetchPageResult, error) {
	page := f.browser.MustPage(url).MustWaitLoad()
	result := &FetchPageResult{
		RodPage: page,
	}
	return result, nil
}

func BuildRodBrowser(browserPath string) (*rod.Browser, error) {
	u := launcher.New().Bin(browserPath).NoSandbox(true).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	fmt.Printf("Start browser: %s\n", u)
	return browser, nil
}

func BuildBrowser(browserPath string) (browser *rod.Browser, cleanup func() error, err error) {
	fmt.Println("get launcher")
	l := launcher.New().
		Bin(browserPath).
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu", "").
		Set("disable-software-rasterizer", "").
		Set("single-process", "").
		Set("homedir", "/tmp").
		Set("data-path", "/tmp/data-path").
		Set("disk-cache-dir", "/tmp/cache-dir")

	launchArgs := l.FormatArgs()
	fmt.Printf("launchArgs: %s\n", launchArgs)

	fmt.Println("start launcher")
	url, err := l.Launch()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to launch browser: %w", err)
	}

	fmt.Printf("url: %s\n", url)
	browser = rod.New().ControlURL(url)
	// .Trace(true)
	fmt.Println("start connect rod")
	if err := browser.Connect(); err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to browser: %w", err)
	}
	fmt.Println("connected rod")

	cleanup = func() error {
		if err := browser.Close(); err != nil {
			return fmt.Errorf("Failed to close browser: %w", err)
		}
		return nil
	}

	return browser, cleanup, nil
}
