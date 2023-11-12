package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector-crawler/pkg/scrapper"
	"github.com/shoet/trends-collector-crawler/pkg/webcrawler"
	"github.com/shoet/trends-collector/config"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/slack"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util/timeutil"
)

func exitFatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		exitFatal(err)
	}

	db := dynamodb.NewFromConfig(cfg)
	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		fmt.Println("failed NewReadClocker")
		exitFatal(err)
	}
	repo := store.NewPageRepository(db, clocker)

	envConfig, err := config.NewConfig()
	if err != nil {
		exitFatal(err)
	}

	httpclient := http.Client{}
	slackClient, err := slack.NewSlackClient(
		&httpclient,
		envConfig.SlackBOTToken,
		envConfig.SlackChannel,
	)
	if err != nil {
		exitFatal(err)
	}

	scrappers, err := buildScrappers(clocker, slackClient)
	if err != nil {
		exitFatal(err)
	}

	c, err := webcrawler.NewWebCrawler(
		&httpclient, envConfig.BrowserPath, scrappers, db, repo)
	if err != nil {
		exitFatal(err)
	}

	err = c.CrawlPages(ctx)
	if err != nil {
		fmt.Println("failed CrawlPages")
		exitFatal(err)
	}

	// TODO: page summary

}

func buildScrappers(
	clocker interfaces.Clocker, slackClient *slack.SlackClient,
) (scrapper.Scrappers, error) {
	dailyTrendsScrapper := scrapper.NewGoogleTrendsDailyTrendsScrapper(clocker)
	realTimeTrendsScrapper := scrapper.NewGoogleTrendsRealTimeTrendsScrapper(clocker)
	hhkbScrapper := scrapper.NewHHKBStudioNotifyScrapper(slackClient)
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
		{
			Category: "HHKB",
			Url:      "https://www.pfu.ricoh.com/direct/hhkb/hhkb-studio/detail_pd-id120b.html",
			Scrapper: hhkbScrapper,
		},
	}
	return scrappers, nil
}
