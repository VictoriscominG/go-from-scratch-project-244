package formatters

import (
	mydiff "code/internal/diff"
	"encoding/json"
)

func Json(result *mydiff.DiffResult) (string, error) {
	got, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(got), nil
}
