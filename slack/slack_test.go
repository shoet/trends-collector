package slack

import (
	"net/http"
	"testing"

	"github.com/shoet/trends-collector/config"
)

func Test_SendMessage(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed load config: %v", err)
	}
	httpClient := http.Client{}
	slackClient, err := NewSlackClient(&httpClient, cfg.SlackBOTToken, cfg.SlackChannel)
	if err != nil {
		t.Fatalf("failed create slack client: %v", err)
	}
	if err := slackClient.SendMessage("test"); err != nil {
		t.Fatalf("failed send message: %v", err)
	}
}
