package scrapper

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/shoet/trends-collector/entities"
	"github.com/shoet/trends-collector/interfaces"
	"github.com/shoet/trends-collector/slack"
	"github.com/shoet/trends-collector/util/timeutil"
)

type ScrapperInput struct {
	RodPage *rod.Page
}

type Scrapper interface {
	ScrapePage(string, *ScrapperInput) ([]*entities.Page, error)
}

type Scrappers []struct {
	Category string
	Url      string
	Scrapper Scrapper
	Pages    []*entities.Page
}

type GoogleTrendsDailyTrendsScrapper struct {
	clocker interfaces.Clocker
}

func NewGoogleTrendsDailyTrendsScrapper(
	clocker interfaces.Clocker,
) *GoogleTrendsDailyTrendsScrapper {
	return &GoogleTrendsDailyTrendsScrapper{clocker: clocker}
}

func (g *GoogleTrendsDailyTrendsScrapper) ScrapePage(
	category string,
	input *ScrapperInput,
) ([]*entities.Page, error) {
	topElement := input.RodPage.MustElement(".feed-list-wrapper")
	trends := topElement.MustElements(".md-list-block")
	pages := make([]*entities.Page, len(trends), len(trends))
	ymd := timeutil.NowFormatYYYYMMDD(g.clocker)
	for i, element := range trends {
		rankStr := element.MustElement(".index").MustText()
		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert rank: %w", err)
		}
		title := element.MustElement(".title").MustText()
		count := element.MustElement(".search-count-title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")

		pages[i] = &entities.Page{
			PartitionKey: ymd,
			TrendRank:    int64(rank),
			Category:     category,
			PageUrl:      *url,
			Title:        title,
			Trend:        count,
		}
	}
	return pages, nil
}

type GoogleTrendsRealTimeTrendsScrapper struct {
	clocker interfaces.Clocker
}

func NewGoogleTrendsRealTimeTrendsScrapper(
	clocker interfaces.Clocker,
) *GoogleTrendsRealTimeTrendsScrapper {
	return &GoogleTrendsRealTimeTrendsScrapper{clocker: clocker}
}

func (g *GoogleTrendsRealTimeTrendsScrapper) ScrapePage(
	category string,
	input *ScrapperInput,
) ([]*entities.Page, error) {
	topElement := input.RodPage.MustElement(".trending-feed-content")
	trends := topElement.MustElements(".feed-item-header")
	pages := make([]*entities.Page, len(trends), len(trends))
	ymdhms := timeutil.NowFormatYYYYMMDDHHMMSS(g.clocker)
	for i, element := range trends {
		rankStr := element.MustElement(".index").MustText()
		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert rank: %w", err)
		}
		title := element.MustElement(".title").MustText()
		url := element.MustElement(".summary-text>a").MustAttribute("href")
		pages[i] = &entities.Page{
			PartitionKey: ymdhms,
			TrendRank:    int64(rank),
			Category:     category,
			PageUrl:      *url,
			Title:        title,
		}
	}
	return pages, nil
}

// TODO: temporary HHKB
type HHKBStudioNotifyScrapper struct {
	slackClient *slack.SlackClient
}

func NewHHKBStudioNotifyScrapper(slackClient *slack.SlackClient) *HHKBStudioNotifyScrapper {
	return &HHKBStudioNotifyScrapper{
		slackClient: slackClient,
	}
}

func (h *HHKBStudioNotifyScrapper) ScrapePage(
	category string,
	input *ScrapperInput,
) ([]*entities.Page, error) {
	// get screenshot
	imageBytes, err := input.RodPage.Screenshot(true, &proto.PageCaptureScreenshot{
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Height: 1000,
			Width:  1500,
			Scale:  1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot: %w", err)
	}
	savePath := filepath.Join(os.TempDir(), "hhkb.png")
	if err := SavePingFile(imageBytes, savePath); err != nil {
		return nil, fmt.Errorf("failed to save screenshot: %w", err)
	}

	html := input.RodPage.MustHTML()
	searchText := "大好評につき現在在庫切れです。"
	message := searchText
	if !strings.Contains(html, searchText) {
		message = "<!channel> HHKBの在庫があります！"
	}
	message = message + "\n" + "https://www.pfu.ricoh.com/direct/hhkb/hhkb-studio/detail_pd-id120b.html"
	if err := h.slackClient.SendMessage(message); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	pages := make([]*entities.Page, 0, 1)
	return pages, nil
}

func SavePingFile(data []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	i, _, _ := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	if err := png.Encode(f, i); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}
