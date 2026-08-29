package formatters

import (
	"bytes"
	mydiff "code/internal/diff"
	"fmt"
	"strings"
)

func Stylish(result *mydiff.DiffResult) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	blockStylish(&buf, result.Items, 1, false)
	buf.WriteString("}")
	return buf.String()
}

// formatBlock рисует список элементов на заданном уровне.
// Параметр isNested означает, что мы находимся внутри мапы, добавленной или
// удалённой целиком. Внутри такого блока вложенные элементы выводятся без
// знаков "+"/"-", но со своими значениями.
func blockStylish(buf *bytes.Buffer, items []mydiff.DiffItem, level int, isNested bool) {
	for _, item := range items {
		// Если есть вложенная структура (мапа), рисуем блок со скобками
		if item.Nested != nil {
			var prefix string
			if isNested {
				prefix = " "
			} else {
				prefix = getPrefix(item.Status)
			}
			indentStr := calcIndent(level)

			buf.WriteString(fmt.Sprintf("%s%s %s: {\n", indentStr, prefix, item.Key))

			// Внутри добавленного/удалённого блока знаки подавляем
			childNested := item.Status == mydiff.ChangeAdded || item.Status == mydiff.ChangeRemoved
			blockStylish(buf, item.Nested.Items, level+1, childNested)

			// Закрываем скобку
			buf.WriteString(fmt.Sprintf("%s  }\n", indentStr))
			continue
		}

		// Обычные строки (не мапы)
		prefix := getPrefix(item.Status)
		indentStr := calcIndent(level)

		// Внутри добавленного/удалённого блока знак не рисуем,
		// но значение (After для added, Before для removed) выводим.
		if isNested {
			valStr := formatValue(item.After)
			if item.Status == mydiff.ChangeRemoved || item.Status == mydiff.ChangeUnchanged {
				valStr = formatValue(item.Before)
			}

			buf.WriteString(fmt.Sprintf("%s  %s: %s\n", indentStr, item.Key, valStr))
			continue
		}

		switch item.Status {
		case mydiff.ChangeUpdated:
			beforeStr := formatValue(item.Before)
			afterStr := formatValue(item.After)

			indentStr := calcIndent(level)

			buf.WriteString(fmt.Sprintf("%s- %s: %s\n", indentStr, item.Key, beforeStr))
			buf.WriteString(fmt.Sprintf("%s+ %s: %s\n", indentStr, item.Key, afterStr))
		default:
			valStr := formatValue(item.After) // Для Added
			if item.Status == mydiff.ChangeRemoved || item.Status == mydiff.ChangeUnchanged {
				valStr = formatValue(item.Before) // Для Removed/Unchanged
			}
			buf.WriteString(fmt.Sprintf("%s%s %s: %s\n", indentStr, prefix, item.Key, valStr))
		}
	}
}

// calcIndent вычисляет строку пробелов по формуле: level * 4 - 2
func calcIndent(level int) string {
	spacesCount := level*4 - 2
	if spacesCount < 0 {
		spacesCount = 0
	}
	return strings.Repeat(" ", spacesCount)
}

// getPrefix по константе возвращает необходимый префикс
func getPrefix(t string) string {
	switch t {
	case mydiff.ChangeAdded:
		return "+"
	case mydiff.ChangeRemoved:
		return "-"
	default:
		return " " // для unchanged
	}
}

// Если значение nil, formatValue возвращает строку null
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}
	// Для простых типов используем стандартное форматирование
	return fmt.Sprintf("%v", v)
}
