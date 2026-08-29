package formatters

import (
	mydiff "code/internal/diff"
	"fmt"
)

func Formatters(result *mydiff.DiffResult, format string) (string, error) {
	var str string
	var err error
	switch format {
	case "stylish":
		str = Stylish(result)
	case "plain":
		str = Plain(result)
	case "json":
		str, err = Json(result)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("the missing output format is specified: %s. Specify the plain, stylish, or json format.", format)
	}
	return str, nil
}
