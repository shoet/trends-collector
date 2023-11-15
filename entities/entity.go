package entities

import "github.com/aws/aws-lambda-go/events"

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type Topic struct {
	Name      string `json:"name" dynamodbav:"name"`
	CreatedAt string `json:"createdAt,omitempty" dynamodbav:"created_at"`
	UpdatedAt string `json:"updatedAt,omitempty" dynamodbav:"updated_at"`
}

type PageId int64

type Page struct {
	PartitionKey string `json:"partition_key" dynamodbav:"partition_key"`
	TrendRank    int64  `json:"trendRank" dynamodbav:"trend_rank"`
	Category     string `json:"category" dynamodbav:"category"`
	Title        string `json:"title" dynamodbav:"title"`
	Trend        string `json:"trend" dynamodbav:"trend"`
	Url          string `json:"url" dynamodbav:"url"`
	CreatedAt    string `json:"createdAt,omitempty" dynamodbav:"created_at"`
	UpdatedAt    string `json:"updatedAt,omitempty" dynamodbav:"updated_at"`
}
