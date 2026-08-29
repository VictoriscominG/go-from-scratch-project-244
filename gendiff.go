package code

import (
	mydiff "code/internal/diff"
	myformatter "code/internal/formatters"
	myparser "code/internal/parser"
)

func GenDiff(path1, path2, format string) (string, error) {
	config1, err := myparser.ParseFile(path1)
	if err != nil {
		return "", err
	}
	config2, err := myparser.ParseFile(path2)
	if err != nil {
		return "", err
	}
	diffResult, err := mydiff.DiffFile(config1, config2)
	if err != nil {
		return "", err
	}

	str, err := myformatter.Formatters(diffResult, format)
	if err != nil {
		return "", err
	}
	return str, nil
}
