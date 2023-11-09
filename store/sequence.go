package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var sequenceTableName = "sequence"

type SequenceId int64

func NextSequence(ctx context.Context, db *dynamodb.Client, tableName string) (SequenceId, error) {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(sequenceTableName),
		Key: map[string]types.AttributeValue{
			"tablename": &types.AttributeValueMemberS{
				Value: tableName,
			},
		},
		UpdateExpression: aws.String("set seq = seq + :val"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberN{
				Value: "1",
			},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}
	item, err := db.UpdateItem(
		ctx,
		updateInput,
	)
	if err != nil {
		return 0, fmt.Errorf("failed update sequence: %w", err)
	}
	var seq struct {
		TableName string     `dynamodbav:"tablename"`
		Seq       SequenceId `dynamodbav:"seq"`
	}
	err = attributevalue.UnmarshalMap(item.Attributes, &seq)
	if err != nil {
		return 0, fmt.Errorf("failed unmarshal sequence: %w", err)
	}
	return seq.Seq, nil
}
