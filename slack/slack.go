package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shoet/trends-collector/interfaces"
)

type SlackClient struct {
	client  interfaces.Client
	token   string
	channel string
}

func NewSlackClient(client interfaces.Client, token string, channel string) (*SlackClient, error) {
	return &SlackClient{
		client:  client,
		token:   token,
		channel: channel,
	}, nil
}

type SendMessageInput struct {
	image []byte
}

type SendMessageResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Message struct {
		Text        string `json:"text"`
		Username    string `json:"username"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			Text     string `json:"text"`
			ID       int    `json:"id"`
			Fallback string `json:"fallback"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"message"`
}

func (s *SlackClient) SendMessage(message string, input *SendMessageInput) error {
	post := struct {
		Text    string `json:"text"`
		Channel string `json:"channel"`
	}{
		Text:    message,
		Channel: s.channel,
	}
	b, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}
	req, err := http.NewRequest(
		"POST",
		"https://slack.com/api/chat.postMessage",
		bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	respB, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	var resp SendMessageResponse
	if err := json.NewDecoder(respB.Body).Decode(&resp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	return nil
}
