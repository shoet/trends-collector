package summary

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/go-rod/rod"
	"github.com/google/uuid"
	"github.com/shoet/trends-collector-crawler/pkg/chatgpt"
	"github.com/shoet/trends-collector-crawler/pkg/fetcher"
	"github.com/shoet/trends-collector/entities"
)

type SummaryGenerator struct {
	chatGPTService *chatgpt.ChatGPTService
	pageFetcher    *fetcher.PageFetcher
}

func NewSummaryGenerator(
	fetcher *fetcher.PageFetcher, chatgpt *chatgpt.ChatGPTService,
) (*SummaryGenerator, error) {
	return &SummaryGenerator{
		chatGPTService: chatgpt,
		pageFetcher:    fetcher,
	}, nil
}

func (s *SummaryGenerator) MakeSummary(url string) (*entities.Summary, error) {
	page := s.pageFetcher.FetchPage(url)
	output, err := ScrapBody(page)
	if err != nil {
		return nil, fmt.Errorf("failed to scrap document: %w", err)
	}

	summaryTemplate, err := summaryTemplateBuilder(output)
	if err != nil {
		return nil, fmt.Errorf("failed to build summary template: %w", err)
	}

	summary, err := s.chatGPTService.ChatCompletions(&chatgpt.ChatCompletionsInput{
		Text: summaryTemplate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get summary: %w", err)
	}
	output.Id = entities.SummaryId(uuid.New().String())
	output.PageUrl = url
	output.Summary = summary

	return output, nil
}

func ScrapBody(page *rod.Page) (*entities.Summary, error) {
	titles, err := page.Elements("h1")
	if err != nil {
		return nil, fmt.Errorf("failed to get h1: %w", err)
	}

	titleBuilder := strings.Builder{}
	for _, t := range titles {
		titleBuilder.WriteString(t.MustText())
	}

	paragraphs, err := page.Elements("p")
	if err != nil {
		return nil, fmt.Errorf("failed to get p: %w", err)
	}

	contentBuilder := strings.Builder{}
	for _, p := range paragraphs {
		contentBuilder.WriteString(p.MustText() + "\n")
	}

	return &entities.Summary{
		Title:   titleBuilder.String(),
		Content: contentBuilder.String(),
	}, nil
}

func summaryTemplateBuilder(v any) (string, error) {
	tmpl, err := template.New("summary").Parse(gptRequestSummaryTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, v); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buffer.String(), nil
}

//go:embed summary_template.txt
var gptRequestSummaryTemplate string
