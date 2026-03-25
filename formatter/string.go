package formatter

import (
	"encoding/json"
	"fmt"
)

func ToJSONStr(v any, indent bool) (string, error) {
	var jsonText []byte
	var err error
	if indent {
		jsonText, err = json.MarshalIndent(v, "", "  ")
	} else {
		jsonText, err = json.Marshal(v)
	}

	if err != nil {
		return "", fmt.Errorf("formatter: to json: %w", err)
	}

	return string(jsonText), nil
}
