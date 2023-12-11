package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util/timeutil"
)

type SummaryRepository struct {
	DB      *dynamodb.Client
	Clocker interfaces.Clocker
}

func NewSummaryRepository(db *dynamodb.Client, clocker interfaces.Clocker) *SummaryRepository {
	return &SummaryRepository{
		DB:      db,
		Clocker: clocker,
	}
}

func (s *SummaryRepository) TableName() string {
	return "summaries"
}

type SummaryRepositoryAddSummaryInput struct {
	Id      entities.SummaryId
	PageUrl string
	Title   string
	Content string
	Summary string
}

func (s *SummaryRepository) AddSummary(
	ctx context.Context,
	input *SummaryRepositoryAddSummaryInput,
) (entities.SummaryId, error) {
	newTopic := &entities.Summary{
		Id:        input.Id,
		PageUrl:   input.PageUrl,
		Title:     input.Title,
		Content:   input.Content,
		Summary:   input.Summary,
		CreatedAt: timeutil.NowFormatRFC3339(s.Clocker),
		UpdatedAt: timeutil.NowFormatRFC3339(s.Clocker),
	}
	av, err := attributevalue.MarshalMap(newTopic)
	if err != nil {
		return "", fmt.Errorf("failed MarshalMap summary: %w", err)
	}
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(s.TableName()),
		Item:      av,
	}
	_, err = s.DB.PutItem(ctx, putInput)
	if err != nil {
		err := fmt.Errorf("failed PutItem summary: %w", err)
		fmt.Println(err.Error())
		return "", err
	}
	return input.Id, nil
}
