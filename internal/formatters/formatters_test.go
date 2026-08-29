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
		err                                   bool
	}{
		{
			name:     "test json file diff format stylish",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-stylish.txt",
			format:   "stylish",
			err:      false,
		}, {
			name:     "test yaml file diff format stylish",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-stylish.txt",
			format:   "stylish",
			err:      false,
		}, {
			name:     "test json file diff format plain",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-plain.txt",
			format:   "plain",
			err:      false,
		}, {
			name:     "test yaml file diff format plain",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-plain.txt",
			format:   "plain",
			err:      false,
		}, {
			name:     "test json file diff format json",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-json.txt",
			format:   "json",
			err:      false,
		}, {
			name:     "test yaml file diff format json",
			before:   "recursive-file1.yaml",
			after:    "recursive-file2.yaml",
			expected: "expected-diff-json.txt",
			format:   "json",
			err:      false,
		}, {
			name:     "test json file diff whith missing format",
			before:   "recursive-file1.json",
			after:    "recursive-file2.json",
			expected: "expected-diff-stylish.txt",
			format:   "xml",
			err:      true,
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
			if tt.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			require.Equal(t, string(expected), renderResult)
		})
	}
}
