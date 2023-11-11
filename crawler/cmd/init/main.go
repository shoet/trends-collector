package main

import (
	"context"
	"fmt"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/shoet/trends-collector/store"
)

func main() {
	ctx := context.Background()
	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		log.Fatal(err)
	}

	db := dynamodb.NewFromConfig(cfg)
	if err := store.AddSequenceTable(ctx, db, "pages"); err != nil {
		log.Fatalf("failed add sequence table [pages]: %s", err.Error())
	}
}
