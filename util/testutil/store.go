package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBForTest(t *testing.T, ctx context.Context, region string) (*dynamodb.Client, error) {
	t.Helper()
	c, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %s\n", err.Error())
	}
	client := dynamodb.NewFromConfig(c)
	return client, nil
}
