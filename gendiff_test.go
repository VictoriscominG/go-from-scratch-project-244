package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	t.Run("test diff json file", func(t *testing.T) {
		result, err := GenDiff(
			"testdata/recursive-file1.json",
			"testdata/recursive-file2.json",
			"plain",
		)
		require.NoError(t, err)

		expected, err := os.ReadFile("testdata/expected-diff-plain.txt")
		require.NoError(t, err)
		require.Equal(t, string(expected), result)
	})

	t.Run("test not json/yaml format file in path1", func(t *testing.T) {
		tempDir := t.TempDir()
		path1 := filepath.Join(tempDir, "txt-file.txt")
		require.NoError(t, os.WriteFile(path1, []byte("{}"), 0o644))

		_, err := GenDiff(path1, "testdata/recursive-file2.json", "plain")

		require.Error(t, err)
	})

	t.Run("test not json/yaml format file in path2", func(t *testing.T) {
		tempDir := t.TempDir()
		path2 := filepath.Join(tempDir, "txt-file.txt")
		require.NoError(t, os.WriteFile(path2, []byte("{}"), 0o644))

		_, err := GenDiff("testdata/recursive-file1.json", path2, "plain")

		require.Error(t, err)
	})

	t.Run("test json file diff whith missing format", func(t *testing.T) {
		_, err := GenDiff(
			"testdata/recursive-file1.json",
			"testdata/recursive-file2.json",
			"xml",
		)
		require.Error(t, err)
	})
}
