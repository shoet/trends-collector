package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/util"
)

func Handler(ctx context.Context, request entities.Request) (entities.Response, error) {
	switch request.HTTPMethod {
	case "GET":
		if request.PathParameters["id"] != "" {
			return getTopic(request)
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

func getTopic(request entities.Request) (entities.Response, error) {
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
