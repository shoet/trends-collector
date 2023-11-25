package push

import (
	"context"

	"github.com/shoet/trends-collector/entities"
)

//go:generate go run github.com/matryer/moq -out push_moq.go . PagesFetcher Pusher
type PagesFetcher interface {
	QueryPageByPartitionKey(ctx context.Context, partitionKey string) ([]*entities.Page, error)
	ScanPageByPartitionKeyPrefix(ctx context.Context, prefix string) ([]string, error)
}

type Pusher interface {
	SendMessage(message string) error
}
