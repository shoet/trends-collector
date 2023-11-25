package entities

import (
	"testing"
)

func Test_Page_FormatTemplate(t *testing.T) {
	page := Page{
		TrendRank: 1,
		Title:     "test",
		Url:       "http://example.com/test",
	}

	template := "第{{.TrendRank}}位: {{.Title}}\n{{.Url}}\n"

	want := "第1位: test\nhttp://example.com/test\n"
	got, err := page.FormatTemplate(template)
	if err != nil {
		t.Errorf("failed to format template: %v", err)
	}

	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
