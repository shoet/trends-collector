package main

import (
	"context"
	"sync"
	"time"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/trends-collector/config"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/slack"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util/timeutil"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, request entities.Request) (Response, error) {
	c, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return Response{Message: "succceed"}, err
	}
	client := dynamodb.NewFromConfig(c)
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to create clocker: %w", err)
	}
	repo := store.NewPageRepository(client, clocker)

	// fetch DailyTrends
	ymd := timeutil.NowFormatYYYYMMDD(clocker)
	var pages entities.Pages
	pages, err = repo.QueryPageByPartitionKey(ctx, ymd)
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to query page: %w", err)
	}

	// sort TrendRank asc
	pages.SortTrendsRank(true)

	cfg, err := config.NewConfig()
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.SlackBOTToken == "" {
		return Response{Message: "failed"}, fmt.Errorf("failed to load config: slack bot token is empty")
	}
	if cfg.SlackChannel == "" {
		return Response{Message: "failed"}, fmt.Errorf("failed to load config: slack channel is empty")
	}
	if cfg.SummaryAPIUrl == "" {
		return Response{Message: "failed"}, fmt.Errorf("failed to load config: summary api url is empty")
	}
	if cfg.SummaryAPIKey == "" {
		return Response{Message: "failed"}, fmt.Errorf("failed to load config: summary api key is empty")
	}
	summaryClient := NewSummaryApiClient(cfg.SummaryAPIUrl, cfg.SummaryAPIKey)

	taskIds := make([]string, len(pages))
	for rankAsc, p := range pages {
		taskId, err := summaryClient.RequestSummaryTask(p.PageUrl)
		if err != nil {
			return Response{Message: "failed"}, fmt.Errorf("failed to request summary: %w", err)
		}
		taskIds[rankAsc] = taskId
	}

	// pooling status
	var wg sync.WaitGroup
	type poolingResult struct {
		Rank   int
		TaskId string
		result *SummaryApiResponse
	}
	ch := make(chan poolingResult, len(taskIds))
	for rankAsc, t := range taskIds {
		wg.Add(1)
		go func(taskId string) {
			defer wg.Done()
			result, err := summaryClient.PoolingTaskStatus(taskId)
			if err != nil {
				return
			}
			ch <- poolingResult{Rank: rankAsc, TaskId: taskId, result: result}
		}(t)
	}
	wg.Wait()

	// 待ち合わせ
	response := make([]*SummaryApiResponse, len(taskIds))
	for res := range ch {
		response[res.Rank] = res.result
		if len(response) == len(taskIds) {
			close(ch)
		}
	}

	// post slack
	httpClient := &http.Client{}
	slackClient, err := slack.NewSlackClient(httpClient, cfg.SlackBOTToken, cfg.SlackChannel)
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to create slack client: %w", err)
	}
	for i, res := range response {
		if res.TaskStatus != "complete" {
			continue
		}
		// post slack
		text := fmt.Sprintf("第%d位【%s】\n%s", i+1, res.PageUrl, res.Summary)
		if err := slackClient.SendMessage(text); err != nil {
			return Response{Message: "failed"}, fmt.Errorf("failed to post slack: %w", err)
		}
	}
	return Response{Message: "succceed"}, nil
}

func main() {
	lambda.Start(Handler)
}

func (s *SummaryApiClient) PoolingTaskStatus(taskId string) (*SummaryApiResponse, error) {
	for {
		resp, err := s.RequestSummaryStatus(taskId)
		if err != nil {
			return nil, fmt.Errorf("failed to request summary: %w", err)
		}
		if resp.TaskStatus == "complete" || resp.TaskStatus == "failed" {
			return resp, nil
		}
		time.Sleep(1 * time.Second)
	}
}

type SummaryApiClient struct {
	client    *http.Client
	apiUrl    string
	apiKey    string
	validator *validator.Validate
}

func NewSummaryApiClient(apiUrl string, apiKey string) *SummaryApiClient {
	return &SummaryApiClient{
		client:    &http.Client{},
		apiUrl:    apiUrl,
		apiKey:    apiKey,
		validator: validator.New(),
	}
}

func (s *SummaryApiClient) RequestSummaryTask(url string) (string, error) {
	requestBody := struct {
		Url string `json:"url"`
	}{
		Url: url,
	}
	b, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}
	httpReq, err := http.NewRequest("POST", s.apiUrl+"/task", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to request summary: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to request summary: status code is %d", resp.StatusCode)
	}
	responseBody := struct {
		TaskID string `json:"task_id" validate:"required"`
	}{}
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(respB, &responseBody); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	if responseBody.TaskID == "" {
		return "", fmt.Errorf("failed to request summary: task_id is empty")
	}
	return responseBody.TaskID, nil
}

func (s *SummaryApiClient) RequestSummaryStatus(taskId string) (*SummaryApiResponse, error) {
	requestBody := struct {
		Id string `json:"id"`
	}{
		Id: taskId,
	}
	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	httpReq, err := http.NewRequest("POST", s.apiUrl+"/get-summary", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("x-api-key", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpClient := http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to request get-summary: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request get-summary: status code is %d", resp.StatusCode)
	}
	responseBody := SummaryApiResponse{}
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(respB, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if err := s.validator.Struct(responseBody); err != nil {
		return nil, fmt.Errorf("failed to validate response: %w", err)
	}
	return &responseBody, nil
}

type SummaryApiResponse struct {
	Id         string `json:"id" validate:"required"`
	TaskStatus string `json:"taskStatus" validate:"required"`
	PageUrl    string `json:"pageUrl"`
	Summary    string `json:"summary"`
}
