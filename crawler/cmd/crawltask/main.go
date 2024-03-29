package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector-crawler/pkg/scrapper"
	"github.com/shoet/trends-collector-crawler/pkg/webcrawler"
	"github.com/shoet/trends-collector/config"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util/timeutil"
)

func exitFatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "local" {
		err := Run()
		if err != nil {
			exitFatal(err)
		}
	} else {
		lambda.Start(Run)
	}
}

func Run() error {
	ctx := context.Background()

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return fmt.Errorf("load aws config: %w", err)
	}

	db := dynamodb.NewFromConfig(cfg)
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		fmt.Println("failed NewReadClocker")
		return fmt.Errorf("failed NewReadClocker: %w", err)
	}
	repo := store.NewPageRepository(db, clocker)

	envConfig, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed NewConfig: %w", err)
	}

	httpclient := http.Client{}
	scrappers, err := buildScrappers(clocker)
	if err != nil {
		return fmt.Errorf("failed buildScrappers: %w", err)
	}

	c, err := webcrawler.NewWebCrawler(
		&httpclient, envConfig.BrowserPath, scrappers, db, repo)
	if err != nil {
		return fmt.Errorf("failed NewWebCrawler: %w", err)
	}
	if c.Closer != nil {
		defer c.Closer()
	}

	err = c.CrawlPages(ctx)
	if err != nil {
		fmt.Println("failed CrawlPages")
		return fmt.Errorf("failed CrawlPages: %w", err)
	}

	return nil
}

func buildScrappers(
	clocker interfaces.Clocker,
) (scrapper.Scrappers, error) {
	dailyTrendsScrapper := scrapper.NewGoogleTrendsDailyTrendsScrapper(clocker)
	realTimeTrendsScrapper := scrapper.NewGoogleTrendsRealTimeTrendsScrapper(clocker)
	scrappers := scrapper.Scrappers{
		{
			Category: "DailyTrends",
			Url:      "https://trends.google.co.jp/trends/trendingsearches/daily?geo=JP&hl=ja",
			Scrapper: dailyTrendsScrapper,
		},
		{
			Category: "RealTimeTrends",
			Url:      "https://trends.google.co.jp/trends/trendingsearches/realtime?geo=JP&hl=ja&category=all",
			Scrapper: realTimeTrendsScrapper,
		},
	}
	return scrappers, nil
}
