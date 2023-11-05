package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/util"
)

func Handler(ctx context.Context, request entities.Request) (entities.Response, error) {
	c, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}
	client := dynamodb.NewFromConfig(c)
	repo := NewTopicRepository(client)
	h := &TopicHandler{
		repo,
	}
	switch request.HTTPMethod {
	case "GET":
		return listTopics(request)
	case "POST":
		return createTopic(request)
	case "DELETE":
		return deleteTopic(request)
	case "PUT":
		return updateTopic(request)
	default:
		return entities.Response{
			StatusCode: 405,
			Body:       fmt.Sprintf("Method %s not allowed", request.HTTPMethod),
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
func listTopics(request entities.Request) (entities.Response, error) {
	resp := []entities.Topic{
		{
			Id:   1,
			Name: "test",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}

func createTopic(request entities.Request) (entities.Response, error) {
	resp := entities.Topic{
		Id: 1,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}

func deleteTopic(request entities.Request) (entities.Response, error) {
	resp := entities.Topic{
		Id: 1,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}

func updateTopic(request entities.Request) (entities.Response, error) {
	resp := []entities.Topic{
		{
			Id:   1,
			Name: "test",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}
