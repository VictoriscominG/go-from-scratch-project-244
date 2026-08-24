package formatters

import mydiff "code/internal/diff"

func Formatters(result *mydiff.DiffResult, format string) (string, error) {
	var str string
	var err error
	switch format {
	case "plain":
		str = Plain(result)
	case "json":
		str, err = Json(result)
		if err != nil {
			return "", err
		}
	default:
		str = Stylish(result)
	}
	return str, nil
}
