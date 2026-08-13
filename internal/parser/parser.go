package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

func ParseFile(path string) (Config, error) {
	// Переводим относительный путь в абсолютный
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize the path: %w", err)
	}

	// Получаем метаданные файла
	info, err := os.Lstat(absPath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil, fmt.Errorf("the file or directory does not exist: %s", path)
		case os.IsPermission(err):
			return nil, fmt.Errorf("no access rights to %s", path)
		default:
			return nil, fmt.Errorf("unknown error when accessing %s: %w", path, err)
		}
	}

	// Обработка символической ссылки
	if info.Mode()&os.ModeSymlink != 0 {
		fmt.Fprintf(os.Stderr, "Warning: %s is a symbolic link, resolving to target\n", absPath)
		info, err = os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("broken symbolic link: %w", err)
		}
	}

	// Проверка, что это обычный файл
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("the path is not a regular file: %s (type: %v)", absPath, info.Mode())
	}

	// Читаем файл
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	var config Config

	// Определяем расширение json/yaml
	ext := filepath.Ext(absPath)

	// Парсим файл в зависимости от расширения
	switch ext {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	default:
		return nil, fmt.Errorf("unsupported extension: %q (file: %s)", ext, absPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s as %s: %w", absPath, ext, err)
	}

	return config, nil
}
