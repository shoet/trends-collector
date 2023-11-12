package webcrawler

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-rod/rod"
	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector-crawler/pkg/scrapper"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/store"
)

type WebCrawler struct {
	client    interfaces.Client
	scrappers scrapper.Scrappers
	browser   *rod.Browser
	db        *dynamodb.Client
	repo      *store.PageRepository
}

func NewWebCrawler(
	client interfaces.Client,
	browserPath string,
	scrappers scrapper.Scrappers,
	db *dynamodb.Client,
	repo *store.PageRepository,
) (*WebCrawler, error) {
	browser, err := fetcher.BuildBrowser(browserPath)
	if err != nil {
		return nil, fmt.Errorf("failed build browser: %w", err)
	}

	return &WebCrawler{
		client:    client,
		scrappers: scrappers,
		browser:   browser,
		db:        db,
		repo:      repo,
	}, nil
}

func (w *WebCrawler) CrawlPages(ctx context.Context) error {
	// TODO: goroutine
	for i := 0; i < len(w.scrappers); i++ {
		s := &w.scrappers[i]
		fmt.Println("Crawl: %s", s.Url)
		page := fetcher.FetchPage(w.browser, s.Url)

		elements, err := s.Scrapper.ScrapePage(s.Category, page)
		if err != nil {
			return fmt.Errorf("failed scrape page: %w", err)
		}

		for _, e := range elements {
			// TODO: remove if exists partition
			input := &store.PageRepositoryAddPageInput{
				Category:  s.Category,
				Partition: e.Partition,
				Title:     e.Title,
				Rank:      e.Rank,
				Trend:     e.Trend,
				Url:       e.Url,
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
