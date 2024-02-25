package fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	cp "github.com/otiai10/copy"
	"github.com/playwright-community/playwright-go"
)

type PlaywrightClient struct {
	browser playwright.Browser
}

type PlaywrightClientConfig struct {
	SkipInstallBrowsers     bool
	BrowserLaunchTimeoutSec int
}

func NewPlaywrightClient(
	config *PlaywrightClientConfig,
) (client *PlaywrightClient, closerFunc func() error, err error) {
	browserBaseDir := "/tmp/playwright/browser"
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: config.SkipInstallBrowsers,
		DriverDirectory:     browserBaseDir,
		Browsers:            []string{"chromium"},
		Verbose:             true,
	}
	if err := playwright.Install(runOption); err != nil {
		return nil, nil, fmt.Errorf("could not install playwright: %v", err)
	}
	pw, err := playwright.Run(runOption)
	if err != nil {
		return nil, nil, fmt.Errorf("could not run playwright: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(browserBaseDir, "chromium-*"))
	if err != nil {
		return nil, nil, fmt.Errorf("could not find browser: %v", err)
	}

	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("could not find browser")
	}
	browserPath := filepath.Join(matches[0], "chrome-linux", "chrome")

	var browserLaunchTimeoutSec float64
	if config.BrowserLaunchTimeoutSec != 0 {
		browserLaunchTimeoutSec = float64(config.BrowserLaunchTimeoutSec) * 1000
	}

	if _, err := os.Stat(browserPath); err != nil {
		return nil, nil, fmt.Errorf("could not find browser: %v", browserPath)
	}
	chromiumOptions := playwright.BrowserTypeLaunchOptions{
		Headless:        playwright.Bool(true),
		ExecutablePath:  playwright.String(browserPath),
		Timeout:         playwright.Float(float64(browserLaunchTimeoutSec)),
		ChromiumSandbox: playwright.Bool(false),
		Args: []string{
			// "--no-sandbox",
			"--single-process",
			"--disable-gpu-sandbox",
		},
	}
	if runtime.GOOS != "linux" {
		chromiumOptions = playwright.BrowserTypeLaunchOptions{}
	}
	fmt.Printf("Start launch browser...")
	browser, err := pw.Chromium.Launch(chromiumOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("could not launch browser: %v", err)
	}
	closer := func() error {
		// if err := browser.Close(); err != nil {
		// 	return fmt.Errorf("could not close browser: %v", err)
		// }
		// if err := pw.Stop(); err != nil {
		// 	return fmt.Errorf("could not stop playwright: %v", err)
		// }
		return nil
	}
	return &PlaywrightClient{
		browser: browser,
	}, closer, nil
}

func CopyBrowser() (string, error) {
	src := "/var/playwright/browser/chromium-1091"
	dst := "/tmp/playwright/browser/chromium-1091"

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err := cp.Copy(src, dst); err != nil {
			return "", fmt.Errorf("could not copy browser: %v", err)
		}
	}
	return dst, nil
}

func (p *PlaywrightClient) FetchPage(url string) (*FetchPageResult, error) {
	page, err := p.browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}
	pageGotoOptions := playwright.PageGotoOptions{
		Timeout: playwright.Float(120000), // ページ表示までのタイムアウト: 2分
	}
	_, err = page.Goto(url, pageGotoOptions)
	if err != nil {
		return nil, fmt.Errorf("could not goto page: %v", err)
	}
	result := &FetchPageResult{
		PlaywrightPage: page,
	}
	return result, nil
}
