package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/util/timeutil"
)

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

func (p *PageRepository) TableName() string {
	return "pages"
}

type PageRepositoryAddPageInput struct {
	PartitionKey string
	TrendRank    int64
	Trend        string
	Category     string
	Title        string
	PageUrl      string
}

func (t *PageRepository) AddPage(
	ctx context.Context,
	input *PageRepositoryAddPageInput,
) (entities.PageId, error) {
	id, err := NextSequence(ctx, t.DB, t.TableName())
	if err != nil {
		return 0, fmt.Errorf("failed NextSequence: %w", err)
	}
	newTopic := &entities.Page{
		PartitionKey: input.PartitionKey,
		TrendRank:    input.TrendRank,
		Category:     input.Category,
		Title:        input.Title,
		Trend:        input.Trend,
		PageUrl:      input.PageUrl,
		CreatedAt:    timeutil.NowFormatRFC3339(t.Clocker),
		UpdatedAt:    timeutil.NowFormatRFC3339(t.Clocker),
	}
	av, err := attributevalue.MarshalMap(newTopic)
	if err != nil {
		return 0, fmt.Errorf("failed MarshalMap page: %w", err)
	}
	putInput := &dynamodb.PutItemInput{
		TableName: aws.String(t.TableName()),
		Item:      av,
	}
	_, err = t.DB.PutItem(ctx, putInput)
	if err != nil {
		err := fmt.Errorf("failed PutItem page: %w", err)
		fmt.Println(err.Error())
		return 0, err
	}
	return entities.PageId(id), nil
}

func (t *PageRepository) QueryPageByPartitionKey(
	ctx context.Context,
	partitionKey string,
) ([]*entities.Page, error) {
	// TODO: urlをpage_urlに変える
	selectExpression := []string{"partition_key", "trend_rank", "title", "page_url"}
	queryRes, err := t.queryPageByPartitionKey(ctx, partitionKey, selectExpression)
	if err != nil {
		return nil, fmt.Errorf("failed queryPageByPartitionKey page: %w", err)
	}

	if queryRes.Count == 0 {
		return []*entities.Page{}, nil
	}

	pages := make([]*entities.Page, len(queryRes.Items), len(queryRes.Items))
	for i, item := range queryRes.Items {
		page := &entities.Page{}
		err = attributevalue.UnmarshalMap(item, page)
		if err != nil {
			return nil, fmt.Errorf("failed UnmarshalMap page: %w", err)
		}
		pages[i] = page
	}
	return pages, nil
}

func (t *PageRepository) queryPageByPartitionKey(
	ctx context.Context,
	partitionKey string,
	selectAttributes []string,
) (*dynamodb.QueryOutput, error) {
	selectExpression := strings.Join(selectAttributes, ",")
	output, err := t.DB.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(t.TableName()),
		KeyConditions: map[string]types.Condition{
			"partition_key": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{
						Value: partitionKey,
					},
				},
			},
		},
		ProjectionExpression: aws.String(selectExpression),
	})
	if err != nil {
		return nil, fmt.Errorf("failed Query: %w", err)
	}
	return output, nil
}

func (t *PageRepository) DeletePageByPartitionKey(
	ctx context.Context,
	partitionKey string,
) error {
	selectExpression := []string{"partition_key", "trend_rank"}
	queryRes, err := t.queryPageByPartitionKey(ctx, partitionKey, selectExpression)
	if err != nil {
		return fmt.Errorf("failed queryPageByPartitionKey page: %w", err)
	}

	if queryRes.Count == 0 {
		return nil
	}

	wr := make([]types.WriteRequest, len(queryRes.Items), len(queryRes.Items))
	for i, item := range queryRes.Items {
		wr[i] = types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{Key: item}}
	}
	_, err = t.DB.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			t.TableName(): wr,
		},
	})
	if err != nil {
		return fmt.Errorf("failed BatchWriteItem page: %w", err)
	}
	return nil
}
