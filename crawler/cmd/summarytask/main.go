package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector-crawler/pkg/chatgpt"
	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector-crawler/pkg/summary"
	"github.com/shoet/trends-collector/config"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/store"
	"github.com/shoet/trends-collector/util/timeutil"
)

func main() {
	ctx := context.Background()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		exitFatal(err)
	}

	db := dynamodb.NewFromConfig(awsCfg)

	clocker, err := timeutil.NewRealClocker()
	if err != nil {
		fmt.Println("failed NewReadClocker")
		exitFatal(err)
	}
	pageRepo := store.NewPageRepository(db, clocker)
	summaryRepo := store.NewSummaryRepository(db, clocker)

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("failed to load config")
		exitFatal(err)
	}
	client := &http.Client{}

	chatgpt := chatgpt.NewChatGPTService(cfg.OpenAIAPIKey, client)
	fetcher, err := fetcher.NewPageFetcher(&fetcher.PageFetcherInput{BrowserPath: cfg.BrowserPath})
	summarizer, err := summary.NewSummaryGenerator(fetcher, chatgpt)
	if err != nil {
		fmt.Println("failed to create summarizer")
		exitFatal(err)
	}

	ymd := timeutil.NowFormatYYYYMMDD(clocker)
	pages, err := pageRepo.QueryPageByPartitionKey(ctx, ymd)
	if err != nil {
		fmt.Println("failed to query page")
		exitFatal(err)
	}

	var wg sync.WaitGroup
	for _, p := range pages {
		wg.Add(1)
		go func(page *entities.Page) {
			output, err := summarizer.MakeSummary(page.PageUrl)
			defer wg.Done()
			if err != nil {
				fmt.Printf("failed to make summary [%s]: %s\n", page.PageUrl, err.Error())
				return
			}
			_, err = summaryRepo.AddSummary(ctx, &store.SummaryRepositoryAddSummaryInput{
				Id:      output.Id,
				PageUrl: output.PageUrl,
				Title:   output.Title,
				Content: output.Content,
				Summary: output.Summary,
			})
			if err != nil {
				fmt.Printf("failed to add summary [%s]: %s\n", page.PageUrl, err.Error())
				return
			}
			// write pages table TODO
			sumaryId := output.Id
			_ = sumaryId
		}(p)
	}

	wg.Wait()
	fmt.Println("Done")
}

func exitFatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}
