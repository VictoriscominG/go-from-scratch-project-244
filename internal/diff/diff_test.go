package diff

import (
	myparser "code/internal/parser"
	fix "code/internal/testfixtures"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiffFile(t *testing.T) {
	// Указываем путь к директории testdata
	basePath := filepath.Join("..", "..", "testdata")

	t.Run("test recursive diff json", func(t *testing.T) {
		beforePath := filepath.Join(basePath, fix.RecursiveFile1JSON)
		afterPath := filepath.Join(basePath, fix.RecursiveFile2JSON)
		expectedPath := filepath.Join(basePath, fix.ExpectedDiffJSON)

		before, err := myparser.ParseFile(beforePath)
		require.NoError(t, err)
		after, err := myparser.ParseFile(afterPath)
		require.NoError(t, err)

		_, err = DiffFile(before, nil)
		require.Error(t, err)
		_, err = DiffFile(nil, after)
		require.Error(t, err)

		got, err := DiffFile(before, after)
		require.NoError(t, err)

		outputJSON, err := json.MarshalIndent(got, "", "  ")
		require.NoError(t, err)

		expectedBytes, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		require.Equal(t, string(expectedBytes), string(outputJSON))
	})

	t.Run("test the scalar has turned into a map", func(t *testing.T) {
		expectedPath := filepath.Join(basePath, fix.ExpectedDiffTest2)
		before := map[string]interface{}{
			"group1": "baz",
			"group2": map[string]interface{}{"group3": "stars"},
			"group3": float64(7),
		}

		after := map[string]interface{}{
			"group1": "baz",
			"group2": map[string]interface{}{"group3": map[string]interface{}{"group4": float64(4)}},
			"group3": float64(7),
		}
		got, err := DiffFile(before, after)
		require.NoError(t, err)

		outputJSON, err := json.MarshalIndent(got, "", "  ")
		require.NoError(t, err)

		expectedBytes, err := os.ReadFile(expectedPath)
		require.NoError(t, err)

		require.Equal(t, string(expectedBytes), string(outputJSON))
	})
}
