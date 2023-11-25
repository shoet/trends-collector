package push

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util/timeutil"
)

type RealTimeTrendsPush struct {
	fetchClient PagesFetcher
	pushClient  Pusher
	clocker     interfaces.Clocker
}

func NewRealTimeTrendsPush(
	fetcher PagesFetcher, client Pusher, clocker interfaces.Clocker) (*RealTimeTrendsPush, error) {
	return &RealTimeTrendsPush{
		fetchClient: fetcher,
		pushClient:  client,
		clocker:     clocker,
	}, nil
}

func (d *RealTimeTrendsPush) Push(ctx context.Context) error {
	pages, err := d.FetchRealTimeTrends(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch daily trends: %w", err)
	}

	message, err := d.formatTrends(pages)
	if err != nil {
		return fmt.Errorf("failed to format trends: %w", err)
	}

	if err := d.pushClient.SendMessage(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (d *RealTimeTrendsPush) FetchRealTimeTrends(ctx context.Context) ([]*entities.Page, error) {
	partitionKey := timeutil.NowFormatYYYYMMDD(d.clocker)
	partitionList, err := d.fetchClient.ScanPageByPartitionKeyPrefix(ctx, partitionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to scan page by partition key prefix: %w", err)
	}
	sort.Slice(partitionList, func(i, j int) bool {
		return partitionList[i] > partitionList[j]
	})
	latest := partitionList[0]
	pages, err := d.fetchClient.QueryPageByPartitionKey(ctx, latest)
	if err != nil {
		return nil, fmt.Errorf("failed to query page by partition key: %w", err)
	}
	return pages, nil
}

func (d *RealTimeTrendsPush) formatTrends(pages []*entities.Page) (string, error) {
	textList := make([]string, len(pages), len(pages))
	template := "第{{.TrendRank}}位: {{.Title}}\n{{.PageUrl}}\n"
	for i, p := range pages {
		f, err := p.FormatTemplate(template)
		if err != nil {
			return "", fmt.Errorf("failed to format template: %w", err)
		}
		textList[i] = f
	}

	header := "只今のリアルタイムトレンドをお届けします。。\n"
	textList = append([]string{header}, textList...)
	message := strings.Join(textList, "\n")
	return message, nil
}
