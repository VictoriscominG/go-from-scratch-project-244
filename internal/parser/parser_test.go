package parser

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	basePath := filepath.Join("..", "..", "testdata")

	t.Run("test recursive file json parse", func(t *testing.T) {
		path := filepath.Join(basePath, "recursive-file1.json")

		cfg, err := ParseFile(path)
		require.NoError(t, err)

		assert.NotNil(t, cfg)

		nested, exists := cfg["group1"]
		assert.True(t, exists)

		child, _ := nested.(map[string]interface{})
		_, exists = child["baz"]
		assert.True(t, exists)
	})

	t.Run("test recursive yaml parse", func(t *testing.T) {
		path := filepath.Join(basePath, "recursive-file1.yaml")

		cfg, err := ParseFile(path)
		require.NoError(t, err)

		assert.NotNil(t, cfg)

		_, exists := cfg["group1"]
		assert.True(t, exists)
	})

	t.Run("test empty path", func(t *testing.T) {
		path := filepath.Join(basePath, "missing-file.json")

		_, err := ParseFile(path)

		require.Error(t, err)
		assert.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("test file with access denide", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "closed.json")
		require.NoError(t, os.WriteFile(path, []byte("{}"), 0o644))
		require.NoError(t, os.Chmod(path, 0o000))

		_, err := ParseFile(path)

		require.Error(t, err)
		assert.True(t, errors.Is(err, os.ErrPermission))
	})

	t.Run("test path is directory", func(t *testing.T) {
		tempDir := t.TempDir()
		dirPath := filepath.Join(tempDir, "test-dir")
		require.NoError(t, os.Mkdir(dirPath, 0o755))

		_, err := ParseFile(dirPath)

		require.Error(t, err)
	})

	t.Run("test not json/yaml format", func(t *testing.T) {
		tempDir := t.TempDir()
		path := filepath.Join(tempDir, "txt-file.txt")
		require.NoError(t, os.WriteFile(path, []byte("{}"), 0o644))

		_, err := ParseFile(path)

		require.Error(t, err)
	})

	t.Run("test func normalizeYamlToJSONTypes", func(t *testing.T) {
		testMap := map[string]interface{}{
			"group1": "baz",
			"group2": map[string]interface{}{"group3": []interface{}{int(4), int64(3)}},
			"group3": int64(7),
		}

		expectedMap := map[string]interface{}{
			"group1": "baz",
			"group2": map[string]interface{}{"group3": []interface{}{float64(4), float64(3)}},
			"group3": float64(7),
		}

		normalized := normalizeYamlToJSONTypes(testMap)
		result, ok := normalized.(map[string]interface{})
		assert.True(t, ok)
		require.True(t, reflect.DeepEqual(expectedMap, result))
	})
}
