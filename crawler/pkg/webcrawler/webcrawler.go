package webcrawler

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector-crawler/pkg/scrapper"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/store"
)

type PageFetcher interface {
	FetchPage(url string) (*fetcher.FetchPageResult, error)
}

type WebCrawler struct {
	client    interfaces.Client
	scrappers scrapper.Scrappers
	fetcher   PageFetcher
	db        *dynamodb.Client
	repo      *store.PageRepository
	Closer    func() error
}

func NewWebCrawler(
	client interfaces.Client,
	browserPath string,
	scrappers scrapper.Scrappers,
	db *dynamodb.Client,
	repo *store.PageRepository,
) (*WebCrawler, error) {
	rodFetcher, cleanup, err := fetcher.NewRodPageFetcher(&fetcher.PageFetcherInput{
		BrowserPath: browserPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build fetcher: %w", err)
	}
	return &WebCrawler{
		client:    client,
		scrappers: scrappers,
		fetcher:   rodFetcher,
		db:        db,
		repo:      repo,
		Closer:    cleanup,
	}, nil
}

func (w *WebCrawler) CrawlPages(ctx context.Context) error {
	// TODO: goroutine
	fmt.Println("Start crawl pages")
	for i := 0; i < len(w.scrappers); i++ {
		s := &w.scrappers[i]
		fmt.Printf("Crawl: %s\n", s.Url)
		result, err := w.fetcher.FetchPage(s.Url)
		if err != nil {
			return fmt.Errorf("failed fetch page: %w", err)
		}

		scrapperInput := &scrapper.ScrapperInput{
			RodPage:        result.RodPage,
			PlaywrightPage: result.PlaywrightPage,
		}

		fmt.Println("Start scrape")
		elements, err := s.Scrapper.ScrapePage(s.Category, scrapperInput)
		if err != nil {
			return fmt.Errorf("failed scrape page: %w", err)
		}

		if len(elements) == 0 {
			fmt.Println("no elements")
			return nil
		}

		fmt.Println("Add to repository")
		if err := w.repo.DeletePageByPartitionKey(ctx, elements[0].PartitionKey); err != nil {
			return fmt.Errorf("failed DeletePageByPartitionKey: %w", err)
		}

		for _, e := range elements {
			input := &store.PageRepositoryAddPageInput{
				PartitionKey: e.PartitionKey,
				TrendRank:    e.TrendRank,
				Category:     s.Category,
				Title:        e.Title,
				Trend:        e.Trend,
				PageUrl:      e.PageUrl,
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
