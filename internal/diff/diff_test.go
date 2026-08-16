package diff

import (
	myparser "code/internal/parser"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffFile(t *testing.T) {
	tests := []struct {
		name     string
		Path1    string
		Path2    string
		expected string
	}{
		{
			name:     "Test diff json",
			Path1:    "file1.json",
			Path2:    "file2.json",
			expected: "expected-diff-json.txt",
		}, {
			name:     "Test diff yaml",
			Path1:    "file1.yaml",
			Path2:    "file2.yaml",
			expected: "expected-diff-yaml.txt",
		},
	}

	// Указываем путь к директории testdata
	basePath := filepath.Join("..", "..", "testdata")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforePath := filepath.Join(basePath, tt.Path1)
			afterPath := filepath.Join(basePath, tt.Path2)
			expectedPath := filepath.Join(basePath, tt.expected)

			before, err := myparser.ParseFile(beforePath)
			require.NoError(t, err)
			after, err := myparser.ParseFile(afterPath)
			require.NoError(t, err)

			_, err = DiffFile(nil, after)
			require.Error(t, err)
			_, err = DiffFile(before, nil)
			require.Error(t, err)

			got, err := DiffFile(before, after)
			require.NoError(t, err)

			expectedBytes, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			require.Equal(t, string(expectedBytes), got)
		})
	}
}
