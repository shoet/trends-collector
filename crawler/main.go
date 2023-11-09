package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util"
)

func exitFatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		exitFatal(err)
	}

	db := dynamodb.NewFromConfig(cfg)
	clocker, err := util.NewRealClocker()
	if err != nil {
		fmt.Println("failed NewReadClocker")
		exitFatal(err)
	}
	repo := store.NewPageRepository(db, clocker)

	scrappers, err := buildScrappers()
	if err != nil {
		exitFatal(err)
	}
	client := http.Client{}
	browserPath := "/opt/homebrew/bin/chromium" // TODO: env

	c := NewWebCrawler(&client, browserPath, scrappers, db, repo)
	err = c.CrawlPages(ctx)
	if err != nil {
		fmt.Println("failed CrawlPages")
		exitFatal(err)
	}

}

func buildScrappers() (Scrappers, error) {
	scrappers := Scrappers{
		{
			category: "DailyTrends",
			url:      "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja",
			scrapper: &GoogleTrendsDailyTrendsScrapper{},
		},
		{
			category: "RealTimeTrends",
			url:      "https://trends.google.co.jp/trends/trendingsearches/realtime?geo=JP&hl=ja&category=all",
			scrapper: &GoogleTrendsRealTimeTrendsScrapper{},
		},
	}
	return scrappers, nil
}
