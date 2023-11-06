package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util"
)

func Handler(ctx context.Context, request entities.Request) (entities.Response, error) {
	c, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}
	client := dynamodb.NewFromConfig(c)
	clocker := &util.RealClocker{}
	repo := store.NewTopicRepository(client, clocker)
	h := &TopicHandler{
		repo,
	}
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			return h.GetTopic(request)
		}
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

type TopicHandler struct {
	repo *store.TopicRepository
}

func (t *TopicHandler) GetTopic(request entities.Request) (entities.Response, error) {
	id := request.PathParameters["id"]

	topic, err := t.repo.GetTopicByName(context.TODO(), id)
	if err != nil {
		fmt.Printf("get topic by id: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}
	resp := []*entities.Topic{
		topic,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}

func listTopics(request entities.Request) (entities.Response, error) {
	resp := []entities.Topic{
		{
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
		Name: "test",
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}

func deleteTopic(request entities.Request) (entities.Response, error) {
	resp := entities.Topic{
		Name: "test",
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
			Name: "test",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return util.ResponseOK(b, nil), nil
}
