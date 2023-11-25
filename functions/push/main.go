package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector/push"
	"github.com/shoet/trends-collector/slack"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util/timeutil"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, event map[string]string) (Response, error) {
	category := event["category"]
	fmt.Println(category)

	client := &http.Client{}

	token, ok := os.LookupEnv("SLACK_BOT_TOKEN")
	if !ok {
		return Response{Message: "failed"}, fmt.Errorf("failed to get slack bot token")
	}
	channel, ok := os.LookupEnv("SLACK_CHANNEL")
	if !ok {
		return Response{Message: "failed"}, fmt.Errorf("failed to get slack channel")
	}

	slackClient, err := slack.NewSlackClient(client, token, channel)
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to create slack client: %w", err)
	}

	c, err := AwsConfig.LoadDefaultConfig(ctx, AwsConfig.WithRegion("ap-northeast-1"))
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to create aws config: %w", err)
	}
	dbClient := dynamodb.NewFromConfig(c)
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		return Response{Message: "failed"}, fmt.Errorf("failed to create clocker: %w", err)
	}
	repo := store.NewPageRepository(dbClient, clocker)

	switch category {
	case "DailyTrends":
		if err := DailyTrendsPush(ctx, repo, slackClient, clocker); err != nil {
			return Response{Message: "failed"}, fmt.Errorf("failed to push daily trends: %w", err)
		}
	case "RealTimeTrends":
		fmt.Println("RealTimeTrends")
	}

	return Response{Message: "success"}, nil
}

func main() {
	lambda.Start(Handler)
}

func DailyTrendsPush(
	ctx context.Context,
	repo *store.PageRepository,
	slackClient *slack.SlackClient,
	clocker *timeutil.RealClocker,
) error {

	dailyTrends, err := push.NewDailyTrendsPush(repo, slackClient, clocker)
	if err != nil {
		return fmt.Errorf("failed to create daily trends push: %w", err)
	}

	if err := dailyTrends.Push(ctx); err != nil {
		return fmt.Errorf("failed to push daily trends: %w", err)
	}
	return nil
}
