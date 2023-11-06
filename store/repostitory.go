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
		fmt.Printf("get item: %s\n", err.Error())
		return nil, err
	}
	var topic *entities.Topic
	err = attributevalue.UnmarshalMap(item.Item, &topic)
	if err != nil {
		fmt.Printf("unmarshal map: %s\n", err.Error())
		return nil, err
	}
	return topic, nil
}

func (t *TopicRepository) AddTopic(
	ctx context.Context, topicName string,
) (entities.TopicId, error) {
	newTopic := &entities.Topic{
		Name:      topicName,
		CreatedAt: util.NowFormatISO8601(t.Clocker),
		UpdatedAt: util.NowFormatISO8601(t.Clocker),
	}
	av, err := attributevalue.MarshalMap(newTopic)
	if err != nil {
		fmt.Printf("marshal map: %s\n", err.Error())
		return 0, err
	}
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName()),
		Item:      av,
	}
	output, err := t.DB.PutItem(ctx, putInput)
	if err != nil {
		fmt.Printf("put item: %s\n", err.Error())
		return 0, err
	}
	fmt.Println(output)
	return 1, nil
}
