package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/util/testutil"
	"github.com/shoet/trends-collector/util/timeutil"
)

func Test_TopicRepository_GetTopicByName(t *testing.T) {
	ctx := context.Background()
	clocker := &timeutil.FixedClocker{}
	client, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("new dynamodb client: %s\n", err.Error())
	}
	want := &entities.Topic{
		Name:      "test",
		CreatedAt: timeutil.NowFormatRFC3339(clocker),
		UpdatedAt: timeutil.NowFormatRFC3339(clocker),
	}

	sut := NewTopicRepository(client, clocker)
	got, err := sut.GetTopicByName(ctx, "test")
	if err != nil {
		t.Fatalf("failed GetTopicByName: %s\n", err.Error())
	}
	testutil.AssertObject(t, want, got)
}

func Test_TopicRepository_AddTopic(t *testing.T) {
	ctx := context.Background()
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		t.Fatalf("new clocker: %s\n", err.Error())
	}
	client, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("new dynamodb client: %s\n", err.Error())
	}
	sut := NewTopicRepository(client, clocker)
	topicId, err := sut.AddTopic(ctx, "test")
	if err != nil {
		t.Fatalf("failed AddTopic: %s\n", err.Error())
	}
	fmt.Println(topicId)
}

func Test_TopicRepository_ListTopics(t *testing.T) {
	ctx := context.Background()
	clocker := &timeutil.FixedClocker{}
	client, err := testutil.NewDynamoDBForTest(t, ctx, "ap-northeast-1")
	if err != nil {
		t.Fatalf("new dynamodb client: %s\n", err.Error())
	}
	sut := NewTopicRepository(client, clocker)
	got, err := sut.ListTopics(ctx, nil)
	if err != nil {
		t.Fatalf("failed ListTopics: %s\n", err.Error())
	}
	fmt.Println(got)
}
