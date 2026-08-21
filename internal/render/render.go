package render

import (
	"bytes"
	"code/internal/diff"
	"fmt"
	"strings"
)

func RenderStylish(result *diff.DiffResult) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	renderBlock(&buf, result.Items, 1, false)
	buf.WriteString("}")
	return buf.String()
}

// renderBlock рисует список элементов на заданном уровне.
// Параметр isNested означает, что мы находимся внутри мапы, добавленной или
// удалённой целиком. Внутри такого блока вложенные элементы выводятся без
// знаков "+"/"-", но со своими значениями.
func renderBlock(buf *bytes.Buffer, items []diff.DiffItem, level int, isNested bool) {
	for _, item := range items {
		// Если есть вложенная структура (мапа), рисуем блок со скобками
		if item.Nested != nil {
			var prefix string
			if isNested {
				prefix = " "
			} else {
				prefix = getPrefix(item.Type)
			}
			indentStr := calcIndent(level)

			buf.WriteString(fmt.Sprintf("%s%s %s: {\n", indentStr, prefix, item.Key))

			// Внутри добавленного/удалённого блока знаки подавляем
			childNested := item.Type == diff.ChangeAdded || item.Type == diff.ChangeRemoved
			renderBlock(buf, item.Nested.Items, level+1, childNested)

			// Закрываем скобку
			buf.WriteString(fmt.Sprintf("%s  }\n", indentStr))
			continue
		}

		// Обычные строки (не мапы)
		prefix := getPrefix(item.Type)
		indentStr := calcIndent(level)

		// Внутри добавленного/удалённого блока знак не рисуем,
		// но значение (After для added, Before для removed) выводим.
		if isNested {
			valStr := formatValue(item.After)
			if item.Type == diff.ChangeRemoved || item.Type == diff.ChangeUnchanged {
				valStr = formatValue(item.Before)
			}

			buf.WriteString(fmt.Sprintf("%s  %s: %s\n", indentStr, item.Key, valStr))
			continue
		}

		switch item.Type {
		case diff.ChangeUpdated:
			beforeStr := formatValue(item.Before)
			afterStr := formatValue(item.After)

			minusIndent := calcIndent(level)
			plusIndent := calcIndent(level)

			buf.WriteString(fmt.Sprintf("%s- %s: %s\n", minusIndent, item.Key, beforeStr))
			buf.WriteString(fmt.Sprintf("%s+ %s: %s\n", plusIndent, item.Key, afterStr))
		default:
			valStr := formatValue(item.After) // Для Added
			if item.Type == diff.ChangeRemoved || item.Type == diff.ChangeUnchanged {
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
	case diff.ChangeAdded:
		return "+"
	case diff.ChangeRemoved:
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
