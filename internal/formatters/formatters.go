package formatters

import (
	mydiff "code/internal/diff"
)

func Formatters(result *mydiff.DiffResult, format string) string {
	var str string
	switch format {
	case "plain":
		str = Plain(result)
	default:
		str = Stylish(result)
	}
	return str
}
