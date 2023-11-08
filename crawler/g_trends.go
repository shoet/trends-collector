package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
)

type Scrapper interface {
	ScrapePage(*rod.Page) ([]*entities.Page, error)
}

type Scrappers []struct {
	id       string
	url      string
	scrapper Scrapper
	result   any
}

type GoogleTrendsCrawler struct {
	client    interfaces.Client
	scrappers Scrappers
	browser   *rod.Browser
}

func NewGoogleTrendsCrawler(client interfaces.Client, browserPath string) *GoogleTrendsCrawler {
	u := launcher.New().Bin(browserPath).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()

	scrappers := Scrappers{
		{
			id:       "DailyTrends",
			url:      "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja",
			scrapper: &GoogleTrendsDailyTrendsScrapper{},
		},
	}
	return &GoogleTrendsCrawler{
		client:    client,
		scrappers: scrappers,
		browser:   browser,
	}
}

func (g *GoogleTrendsCrawler) FetchPage(url string) (*entities.Page, error) {
	return nil, nil
}

// TODO: not use
func (g *GoogleTrendsCrawler) FetchDocument(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed create request: %w", err)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed request: %d", resp.StatusCode)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func (g *GoogleTrendsCrawler) CrawlPages() (Scrappers, error) {
	for _, s := range g.scrappers {
		page := g.browser.MustPage(s.url)
		page.MustWaitLoad()

		pages, err := s.scrapper.ScrapePage(page)
		if err != nil {
			return nil, fmt.Errorf("failed scrape page: %w", err)
		}
		s.result = pages
	}
	return g.scrappers, nil
}

func (g *GoogleTrendsCrawler) GetPageHTML(url string) (*goquery.Document, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed chromedb.Run: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
