package push

import (
	"context"
	"fmt"
	"strings"

	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util/timeutil"
)

type PagesFetcher interface {
	QueryPageByPartitionKey(ctx context.Context, partitionKey string) ([]*entities.Page, error)
}

type Pusher interface {
	SendMessage(message string) error
}

type DailyTrendsPush struct {
	fetchClient PagesFetcher
	pushClient  Pusher
	clocker     interfaces.Clocker
}

func NewDailyTrendsPush(
	fetcher PagesFetcher, client Pusher, clocker interfaces.Clocker) (*DailyTrendsPush, error) {
	return &DailyTrendsPush{
		fetchClient: fetcher,
		pushClient:  client,
		clocker:     clocker,
	}, nil
}

func (d *DailyTrendsPush) Push(ctx context.Context) error {
	pages, err := d.FetchDailyTrends(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch daily trends: %w", err)
	}

	message, err := formatTrends(pages)
	if err != nil {
		return fmt.Errorf("failed to format trends: %w", err)
	}

	if err := d.pushClient.SendMessage(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (d *DailyTrendsPush) FetchDailyTrends(ctx context.Context) ([]*entities.Page, error) {
	partitionKey := timeutil.NowFormatYYYYMMDD(d.clocker)
	pages, err := d.fetchClient.QueryPageByPartitionKey(ctx, partitionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to query page by partition key: %w", err)
	}
	return pages, nil
}

func formatTrends(pages []*entities.Page) (string, error) {
	textList := make([]string, len(pages), len(pages))
	template := "第{{.TrendRank}}位: {{.Title}}\n{{.PageUrl}}\n"
	for i, p := range pages {
		f, err := p.FormatTemplate(template)
		if err != nil {
			return "", fmt.Errorf("failed to format template: %w", err)
		}
		textList[i] = f
	}

	header := "本日のデイリートレンドをお届けします。\n"
	textList = append([]string{header}, textList...)
	message := strings.Join(textList, "\n")
	return message, nil
}
