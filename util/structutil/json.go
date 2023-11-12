package structutil

import (
	"encoding/json"
	"strings"
)

func JSONStrToStruct(s string, v any) error {
	return json.NewDecoder(strings.NewReader(s)).Decode(v)
}
