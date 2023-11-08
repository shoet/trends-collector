package main

import (
	"fmt"
	"net/http"
	"os"
)

func exitFatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	client := http.Client{}
	browserPath := "/opt/homebrew/bin/chromium"
	c := NewGoogleTrendsCrawler(&client, browserPath)
	scrappers, err := c.CrawlPages()
	if err != nil {
		exitFatal(err)
	}
	for _, s := range scrappers {
		r := s.result
		// TODO: add dynamodb
		fmt.Println(r)
	}
}
