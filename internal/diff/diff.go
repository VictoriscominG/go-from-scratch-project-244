package diff

import (
	"fmt"
	"sort"
	"strings"
)

func DiffFile(conf1, conf2 map[string]interface{}) (string, error) {
	if conf1 == nil {
		return "", fmt.Errorf("I can’t compare an empty configuration: %v", conf1)
	}
	if conf2 == nil {
		return "", fmt.Errorf("I can’t compare an empty configuration: %v", conf2)
	}

	// Собираем все уникальные ключи из обоих конфигов
	keysMap := make(map[string]struct{})
	for k := range conf1 {
		keysMap[k] = struct{}{}
	}
	for k := range conf2 {
		keysMap[k] = struct{}{}
	}

	// Превращаем ключи в срез и сортируем по алфавиту (для детерминированного вывода)
	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Создаём слайс строк для вывода, вносим данные
	var lines []string
	lines = append(lines, "{")

	for _, k := range keys {
		v1, ok1 := conf1[k]
		v2, ok2 := conf2[k]

		switch {
		case !ok1 && ok2:
			// Добавлено во втором конфиге
			lines = append(lines, fmt.Sprintf("  + %s: %v", k, v2))
		case ok1 && !ok2:
			// Удалено во втором конфиге
			lines = append(lines, fmt.Sprintf("  - %s: %v", k, v1))
		default:
			// Есть в обоих — сравниваем значения
			if v1 == v2 {
				lines = append(lines, fmt.Sprintf("    %s: %v", k, v1)) // без изменений
			} else {
				lines = append(lines, fmt.Sprintf("  - %s: %v", k, v1))
				lines = append(lines, fmt.Sprintf("  + %s: %v", k, v2))
			}
		}
	}
	lines = append(lines, "}\n")

	return strings.Join(lines, "\n"), nil
}
