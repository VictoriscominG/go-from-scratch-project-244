package formatters

import (
	mydiff "code/internal/diff"
	myparser "code/internal/parser"
	fix "code/internal/testfixtures"
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
			before:   fix.RecursiveFile1JSON,
			after:    fix.RecursiveFile2JSON,
			expected: fix.ExpectedDiffStylish,
			format:   "stylish",
			err:      false,
		}, {
			name:     "test yaml file diff format stylish",
			before:   fix.RecursiveFile1YAML,
			after:    fix.RecursiveFile2YAML,
			expected: fix.ExpectedDiffStylish,
			format:   "stylish",
			err:      false,
		}, {
			name:     "test json file diff format plain",
			before:   fix.RecursiveFile1JSON,
			after:    fix.RecursiveFile2JSON,
			expected: fix.ExpectedDiffPlain,
			format:   "plain",
			err:      false,
		}, {
			name:     "test yaml file diff format plain",
			before:   fix.RecursiveFile1YAML,
			after:    fix.RecursiveFile2YAML,
			expected: fix.ExpectedDiffPlain,
			format:   "plain",
			err:      false,
		}, {
			name:     "test json file diff format json",
			before:   fix.RecursiveFile1JSON,
			after:    fix.RecursiveFile2JSON,
			expected: fix.ExpectedDiffJSON,
			format:   "json",
			err:      false,
		}, {
			name:     "test yaml file diff format json",
			before:   fix.RecursiveFile1YAML,
			after:    fix.RecursiveFile2YAML,
			expected: fix.ExpectedDiffJSON,
			format:   "json",
			err:      false,
		}, {
			name:     "test json file diff whith missing format",
			before:   fix.RecursiveFile1JSON,
			after:    fix.RecursiveFile2JSON,
			expected: fix.ExpectedDiffStylish,
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
