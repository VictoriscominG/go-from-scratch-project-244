package formatters

import (
	mydiff "code/internal/diff"
	"fmt"
	"strings"
)

func Plain(result *mydiff.DiffResult) string {
	lines := keyItems(result.Items, "")
	return strings.Join(lines, "\n")
}

func keyItems(items []mydiff.DiffItem, prefix string) []string {
	var lines []string
	i := 0
	for i < len(items) {
		if i+1 < len(items) &&
			items[i].Status == mydiff.ChangeRemoved &&
			items[i+1].Status == mydiff.ChangeAdded &&
			items[i].Key == items[i+1].Key {

			fullKey := buildKey(prefix, items[i].Key)
			before := normalizeValue(items[i].Before, items[i].Nested)
			after := normalizeValue(items[i+1].After, items[i+1].Nested)
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s", fullKey, before, after))
			i += 2
			continue
		}

		item := items[i]
		fullKey := buildKey(prefix, item.Key)
		switch item.Status {
		case mydiff.ChangeUnchanged:
		case mydiff.ChangeNested:
			lines = append(lines, keyItems(item.Nested.Items, fullKey)...)
		case mydiff.ChangeAdded:
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s", fullKey, normalizeValue(item.After, item.Nested)))
		case mydiff.ChangeRemoved:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", fullKey))
		case mydiff.ChangeUpdated:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s", fullKey, normalizeValue(item.Before, nil), normalizeValue(item.After, item.Nested)))
		}
		i++
	}
	return lines
}

func buildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func normalizeValue(v interface{}, nested *mydiff.DiffResult) string {
	if v == nil && nested != nil {
		return "[complex value]"
	}
	return normalize(v)
}

func normalize(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}
