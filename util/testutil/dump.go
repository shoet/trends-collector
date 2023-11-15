package testutil

import (
	"fmt"
	"io"

	"github.com/go-rod/rod"
	"github.com/microcosm-cc/bluemonday"
)

func DumpRodElement(w io.Writer, elem interface{}) error {
	switch value := elem.(type) {
	case *rod.Element:
		elements := []*rod.Element{value}
		return dumpRodElements(w, elements)
	case rod.Elements:
		return dumpRodElements(w, value)
	default:
		return fmt.Errorf("unknown type: %T", value)
	}

}

func dumpRodElements(w io.Writer, elem []*rod.Element) error {
	for _, e := range elem {
		html, err := e.HTML()
		if err != nil {
			return fmt.Errorf("failed to get html: %w", err)
		}

		p := bluemonday.UGCPolicy()
		p.AllowAttrs("class", "id").OnElements("div")

		html = p.Sanitize(html)

		if err != nil {
			return fmt.Errorf("failed to get html: %w", err)
		}
		_, err = io.WriteString(w, html)
		if err != nil {
			return fmt.Errorf("failed to write html: %w", err)
		}
	}
	return nil
}
