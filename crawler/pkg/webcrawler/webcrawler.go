package webcrawler

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
		fmt.Println("Crawl: %s", s.Url)
		page := w.browser.MustPage(s.Url)
		page.MustWaitLoad()

		pages, err := s.Scrapper.ScrapePage(s.Category, page)
		if err != nil {
			return fmt.Errorf("failed scrape page: %w", err)
		}

		for _, p := range pages {
			input := &store.PageRepositoryAddPageInput{
				Category:  s.Category,
				Partition: p.Partition,
				Title:     p.Title,
				Rank:      p.Rank,
				Trend:     p.Trend,
				Url:       p.Url,
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
