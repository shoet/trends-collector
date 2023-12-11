package main

import (
	"fmt"
	"testing"

	"github.com/shoet/trends-collector/config"
)

func Test_RequestSummaryTask(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Errorf("Error while loading config: %v", err)
	}
	sut := NewSummaryApiClient(cfg.SummaryAPIUrl, cfg.SummaryAPIKey)
	url := "https://mainichi.jp/articles/20231208/ddl/k23/040/141000c"
	taskId, err := sut.RequestSummaryTask(url)
	if err != nil {
		t.Errorf("Error while requesting summary task: %v", err)
	}
	t.Logf("TaskId: %v", taskId)
}

func Test_RequestSummaryStatus(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Errorf("Error while loading config: %v", err)
	}
	sut := NewSummaryApiClient(cfg.SummaryAPIUrl, cfg.SummaryAPIKey)
	result, err := sut.RequestSummaryStatus("5624e8d9-835a-4af9-9aeb-769bbb5c94de")
	if err != nil {
		t.Errorf("Error while requesting summary status: %v", err)
	}
	fmt.Println(result)
}
