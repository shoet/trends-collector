package push

import (
	"context"
	"testing"

	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/util/timeutil"
)

func Test_NewRealTimeTrendsPush_FetchRealTimeTrends(t *testing.T) {
	fetchClientMock := &PagesFetcherMock{}
	fetchClientMock.ScanPageByPartitionKeyPrefixFunc = func(ctx context.Context, prefix string) ([]string, error) {
		return []string{
			"20231125110121",
			"20231125",
			"20231125093903",
			"20231125100123",
		}, nil
	}
	fetchClientMock.QueryPageByPartitionKeyFunc = func(ctx context.Context, partitionKey string) ([]*entities.Page, error) {
		if partitionKey != "20231125110121" {
			t.Errorf("unexpected partition key: %s", partitionKey)
		}
		return []*entities.Page{}, nil
	}

	clocker := timeutil.FixedClocker{}
	sut := &RealTimeTrendsPush{
		fetchClient: fetchClientMock,
		clocker:     &clocker,
	}
	pages, err := sut.FetchRealTimeTrends(context.Background())
	if err != nil {
		t.Errorf("failed to fetch real time trends: %v", err)
	}
	_ = pages
}
