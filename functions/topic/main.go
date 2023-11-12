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
	"github.com/shoet/trends-collector/util/responseutil"
	"github.com/shoet/trends-collector/util/structutil"
	"github.com/shoet/trends-collector/util/timeutil"
)

func Handler(ctx context.Context, request entities.Request) (entities.Response, error) {
	c, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}
	client := dynamodb.NewFromConfig(c)
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		fmt.Printf("new clocker: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}
	repo := store.NewTopicRepository(client, clocker)
	h := &TopicHandler{
		repo,
	}
	switch request.HTTPMethod {
	case "GET":
		return h.ListTopics(request)
	case "POST":
		return h.CreateTopic(request)
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
	return responseutil.ResponseOK(b, nil), nil
}

func (t *TopicHandler) ListTopics(request entities.Request) (entities.Response, error) {
	topics, err := t.repo.ListTopics(context.TODO(), nil)
	b, err := json.Marshal(topics)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return responseutil.ResponseOK(b, nil), nil
}

func (t *TopicHandler) CreateTopic(request entities.Request) (entities.Response, error) {
	body := struct {
		Name string `json:"name"`
	}{}

	if err := structutil.JSONStrToStruct(request.Body, &body); err != nil {
		fmt.Printf("failed deserialize body: %s\n", err.Error())
		return entities.Response{StatusCode: 500}, err
	}

	name, err := t.repo.AddTopic(context.TODO(), body.Name)
	resp := entities.Topic{Name: name}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return responseutil.ResponseOK(b, nil), nil // TODO: status 201
}

func deleteTopic(request entities.Request) (entities.Response, error) {
	resp := entities.Topic{
		Name: "test",
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return entities.Response{StatusCode: 404}, err
	}
	return responseutil.ResponseOK(b, nil), nil
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
	return responseutil.ResponseOK(b, nil), nil
}
