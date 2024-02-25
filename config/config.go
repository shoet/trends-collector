package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	BrowserPath   string `env:"CRAWLER_BROWSER_PATH" envDefault:"/usr/bin/chromium"`
	SlackBOTToken string `env:"SLACK_BOT_TOKEN"`
	SlackChannel  string `env:"SLACK_CHANNEL"`
	OpenAIAPIKey  string `env:"CRAWLER_OPENAI_API_KEY"`
	SummaryAPIUrl string `env:"WEB_PAGE_SUMMARY_API_URL"`
	SummaryAPIKey string `env:"WEB_PAGE_SUMMARY_API_KEY"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed Parse config: %w", err)
	}
	return cfg, nil
}
