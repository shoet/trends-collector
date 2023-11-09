package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/store"
)

type Scrapper interface {
	ScrapePage(string, *rod.Page) ([]*entities.Page, error)
}

type Scrappers []struct {
	category string
	url      string
	scrapper Scrapper
	pages    []*entities.Page
}

type WebCrawler struct {
	client    interfaces.Client
	scrappers Scrappers
	browser   *rod.Browser
	db        *dynamodb.Client
	repo      *store.PageRepository
}

func NewWebCrawler(
	client interfaces.Client,
	browserPath string,
	scrappers Scrappers,
	db *dynamodb.Client,
	repo *store.PageRepository,
) *WebCrawler {
	u := launcher.New().Bin(browserPath).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	fmt.Printf("Start browser: %s\n", u)

	return &WebCrawler{
		client:    client,
		scrappers: scrappers,
		browser:   browser,
		db:        db,
		repo:      repo,
	}
}

func (w *WebCrawler) CrawlPages(ctx context.Context) error {
	// TODO: goroutine
	for i := 0; i < len(w.scrappers); i++ {
		s := &w.scrappers[i]
		fmt.Println("Crawl: %s", s.url)
		page := w.browser.MustPage(s.url)
		page.MustWaitLoad()

		pages, err := s.scrapper.ScrapePage(s.category, page)
		if err != nil {
			return fmt.Errorf("failed scrape page: %w", err)
		}

		for _, p := range pages {
			input := &store.PageRepositoryAddPageInput{
				Category: s.category,
				Title:    p.Title,
				Trend:    p.Trend,
				Url:      s.url,
			}
			_, err := w.repo.AddPage(ctx, input)
			if err != nil {
				// TODO: out log
				fmt.Println("failed AddPage")
				fmt.Println(err)
			}
		}
	}
	return nil
}
