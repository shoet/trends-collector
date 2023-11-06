package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/testutil"
	"github.com/shoet/trends-collector/util"
)

func Test_TopicRepository_GetTopicByName(t *testing.T) {
	ctx := context.Background()
	clocker := &util.FixedClocker{}
	client, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("new dynamodb client: %s\n", err.Error())
	}
	want := &entities.Topic{
		Name:      "test",
		CreatedAt: util.NowFormatISO8601(clocker),
		UpdatedAt: util.NowFormatISO8601(clocker),
	}

	sut := NewTopicRepository(client, clocker)
	got, err := sut.GetTopicByName(ctx, "test")
	if err != nil {
		t.Fatalf("get topic by id: %s\n", err.Error())
	}
	testutil.AssertObject(t, want, got)
}

func Test_TopicRepository_AddTopic(t *testing.T) {
	ctx := context.Background()
	clocker := &util.FixedClocker{}
	client, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("new dynamodb client: %s\n", err.Error())
	}
	sut := NewTopicRepository(client, clocker)
	topicId, err := sut.AddTopic(ctx, "test")
	if err != nil {
		t.Fatalf("get topic by id: %s\n", err.Error())
	}
	fmt.Println(topicId)
}
