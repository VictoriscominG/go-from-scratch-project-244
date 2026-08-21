package render

import (
	mydiff "code/internal/diff"
	myparser "code/internal/parser"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	basePath := filepath.Join("..", "..", "testdata")
	t.Run("test json diff render format stylish", func(t *testing.T) {
		beforePath := filepath.Join(basePath, "recursive-file1.json")
		afterPath := filepath.Join(basePath, "recursive-file2.json")
		expectedPath := filepath.Join(basePath, "expected-diff-recursive-json.txt")

		before, err := myparser.ParseFile(beforePath)
		require.NoError(t, err)
		after, err := myparser.ParseFile(afterPath)
		require.NoError(t, err)

		diffResult, err := mydiff.DiffFile(before, after)
		require.NoError(t, err)

		renderResult := RenderStylish(diffResult)

		expected, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		require.Equal(t, string(expected), renderResult)
	})
}
