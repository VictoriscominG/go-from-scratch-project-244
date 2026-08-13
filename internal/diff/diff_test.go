package diff

import (
	myparser "code/internal/parser"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffFile(t *testing.T) {
	basePath := filepath.Join("..", "..", "testdata")
	beforePath := filepath.Join(basePath, "file1.json")
	afterPath := filepath.Join(basePath, "file2.json")
	expectedPath := filepath.Join(basePath, "expected-diff.txt")

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
}
