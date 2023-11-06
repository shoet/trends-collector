package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util"
)

type TopicRepository struct {
	DB      *dynamodb.Client
	Clocker interfaces.Clocker
}

func NewTopicRepository(db *dynamodb.Client, clocker interfaces.Clocker) *TopicRepository {
	return &TopicRepository{
		DB:      db,
		Clocker: clocker,
	}
}

func (t *TopicRepository) TableName() string {
	return "Topics"
}

func (t *TopicRepository) GetTopicByName(ctx context.Context, name string) (*entities.Topic, error) {
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(t.TableName()),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{
				Value: name,
			},
		},
	}
	item, err := t.DB.GetItem(
		ctx,
		getInput,
	)
	if err != nil {
		err := fmt.Errorf("failed GetItem topics: %w", err)
		fmt.Println(err.Error())
		return nil, err
	}
	var topic *entities.Topic
	err = attributevalue.UnmarshalMap(item.Item, &topic)
	if err != nil {
		err := fmt.Errorf("failed UnmarshalMap topics: %w", err)
		fmt.Println(err.Error())
		return nil, err
	}
	return topic, nil
}

type ListTopicsInput struct {
	Limit int
}

func (t *TopicRepository) ListTopics(
	ctx context.Context, option *ListTopicsInput,
) ([]*entities.Topic, error) {
	defaultOpt := &ListTopicsInput{
		Limit: 100,
	}
	util.MergeStruct(defaultOpt, option)
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(t.TableName()),
		Limit:     aws.Int32(int32(defaultOpt.Limit)),
	}
	scanOutput, err := t.DB.Scan(ctx, scanInput)
	if err != nil {
		err := fmt.Errorf("failed Scan topics: %w", err)
		fmt.Println(err.Error())
		return nil, err
	}
	var topics []*entities.Topic
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &topics)
	if err != nil {
		err := fmt.Errorf("failed UnmarshalListOfMaps topics: %w", err)
		fmt.Println(err.Error())
		return nil, err
	}
	return topics, nil
}

func (t *TopicRepository) AddTopic(
	ctx context.Context, topicName string,
) (string, error) {
	newTopic := &entities.Topic{
		Name:      topicName,
		CreatedAt: util.NowFormatRFC3339(t.Clocker),
		UpdatedAt: util.NowFormatRFC3339(t.Clocker),
	}
	av, err := attributevalue.MarshalMap(newTopic)
	if err != nil {
		err := fmt.Errorf("failed MarshalMap topic: %w", err)
		fmt.Println(err.Error())
		return "", err
	}
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName()),
		Item:      av,
	}
	_, err = t.DB.PutItem(ctx, putInput)
	if err != nil {
		err := fmt.Errorf("failed PutItem topic: %w", err)
		fmt.Println(err.Error())
		return "", err
	}
	return topicName, nil
}

type PageRepository struct {
	DB      *dynamodb.Client
	Clocker interfaces.Clocker
}

func NewPageRepository(db *dynamodb.Client, clocker interfaces.Clocker) *PageRepository {
	return &PageRepository{
		DB:      db,
		Clocker: clocker,
	}
}
