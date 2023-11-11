package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoet/trends-collector/util/testutil"
)

func Test_NextSequence(t *testing.T) {
	table := "test"
	ctx := context.Background()
	db, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("failed to create dynamodb client: %v", err)
	}
	nextId, err := NextSequence(ctx, db, table)
	if err != nil {
		t.Fatalf("failed to get next sequence: %v", err)
	}
	fmt.Printf("nextId: %v\n", nextId)
}
