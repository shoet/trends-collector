package entities

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
)

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
	PageUrl      string `json:"page_url" dynamodbav:"page_url"`
	CreatedAt    string `json:"createdAt,omitempty" dynamodbav:"created_at"`
	UpdatedAt    string `json:"updatedAt,omitempty" dynamodbav:"updated_at"`
}

func (p *Page) FormatTemplate(templateText string) (string, error) {
	tmpl, err := template.New("page").Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, p)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	return buffer.String(), nil
}

type Pages []*Page

func (p Pages) SortTrendsRank(asc bool) {
	if asc {
		sort.Slice(p, func(i, j int) bool { return p[i].TrendRank < p[j].TrendRank })
	} else {
		sort.Slice(p, func(i, j int) bool { return p[i].TrendRank > p[j].TrendRank })
	}
}

type SummaryId string

type Summary struct {
	Id        SummaryId `json:"id" dynamodbav:"id"`
	PageUrl   string    `json:"pageUrl" dynamodbav:"page_url"`
	Title     string    `json:"title" dynamodbav:"title"`
	Content   string    `json:"content" dynamodbav:"content"`
	Summary   string    `json:"summary" dynamodbav:"summary"`
	CreatedAt string    `json:"createdAt,omitempty" dynamodbav:"created_at"`
	UpdatedAt string    `json:"updatedAt,omitempty" dynamodbav:"updated_at"`
}
