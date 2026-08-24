package formatters

import (
	mydiff "code/internal/diff"
	myparser "code/internal/parser"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatters(t *testing.T) {
	tests := []struct {
		name, before, after, expected, format string
	}{
		{
			name:     "test json file diff format stylish",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-stylish.txt",
			format:   "stylish",
		}, {
			name:     "test yaml file diff format stylish",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-stylish.txt",
			format:   "stylish",
		}, {
			name:     "test json file diff format plain",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-plain.txt",
			format:   "plain",
		}, {
			name:     "test yaml file diff format plain",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-plain.txt",
			format:   "plain",
		}, {
			name:     "test json file diff format json",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-json.txt",
			format:   "json",
		}, {
			name:     "test yaml file diff format json",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-json.txt",
			format:   "json",
		}, {
			name:     "test json file diff whithout format",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-stylish.txt",
			format:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath := filepath.Join("..", "..", "testdata")
			beforePath := filepath.Join(basePath, tt.before)
			afterPath := filepath.Join(basePath, tt.after)
			expectedPath := filepath.Join(basePath, tt.expected)

			before, err := myparser.ParseFile(beforePath)
			require.NoError(t, err)
			after, err := myparser.ParseFile(afterPath)
			require.NoError(t, err)

			diffResult, err := mydiff.DiffFile(before, after)
			require.NoError(t, err)

			renderResult, err := Formatters(diffResult, tt.format)
			require.NoError(t, err)

			expected, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			require.Equal(t, string(expected), renderResult)
		})
	}
}
