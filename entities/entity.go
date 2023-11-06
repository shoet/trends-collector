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
	Id      PageId `json:"id" dynamodbav:"id"`
	Html    string `json:"html" dynamodbav:"html"`
	Summary string `json:"summary" dynamodbav:"summary"`
}

func (p *Page) String() string {
	return "Page"
}

func (p *Page) Tree() map[string]any {
	return map[string]any{}
}
