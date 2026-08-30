package diff

import (
	"fmt"
	"reflect"
	"sort"
)

const (
	ChangeAdded     = "added"
	ChangeRemoved   = "removed"
	ChangeUpdated   = "updated"
	ChangeUnchanged = "unchanged"
	ChangeNested    = "nested"
)

type DiffItem struct {
	Key    string      `json:"key" yaml:"key"`
	Status string      `json:"status" yaml:"status"`
	Before interface{} `json:"before,omitempty" yaml:"before,omitempty"`
	After  interface{} `json:"after,omitempty" yaml:"after,omitempty"`
	Nested *DiffResult `json:"nested,omitempty" yaml:"nested,omitempty"`
}

type DiffResult struct {
	Items []DiffItem `json:"items" yaml:"items"`
}

func DiffFile(config1, config2 map[string]interface{}) (*DiffResult, error) {
	if config1 == nil {
		return nil, fmt.Errorf("I can’t compare an empty configuration: %v", config1)
	}
	if config2 == nil {
		return nil, fmt.Errorf("I can’t compare an empty configuration: %v", config2)
	}

	keys := keysCollect(config1, config2)

	// Объявляем коллекцию []DiffItem для сбора элементов diff
	var items []DiffItem

	for _, k := range keys {
		v1, ok1 := config1[k]
		v2, ok2 := config2[k]

		// Проверяем тип значения на мапу
		map1, isMap1 := v1.(map[string]interface{})
		map2, isMap2 := v2.(map[string]interface{})

		// Объявляем временный объект для одного ключа на текущей итерации
		var item DiffItem
		item.Key = k

		switch {
		case !ok1 && ok2:
			// Добавлено
			item.Status = ChangeAdded
			if isMap2 {
				// Добавили мапу целиком: внутри неё ничего не сравнивается,
				// все вложенные элементы помечаем как added (со значениями).
				nestedResult, err := diffMapAgainstEmpty(map2, false)
				if err != nil {
					return nil, err
				}
				item.Nested = nestedResult
			} else {
				// Добавленный скаляр
				item.After = v2
			}
		case ok1 && !ok2:
			// Удалено
			item.Status = ChangeRemoved
			if isMap1 {
				// Удалили мапу целиком: все вложенные элементы — removed (со значениями).
				nestedResult, err := diffMapAgainstEmpty(map1, true)
				if err != nil {
					return nil, err
				}
				item.Nested = nestedResult
			} else {
				// Удалённый скаляр
				item.Before = v1
			}
		case isMap1 && isMap2:
			// Вложенная структура присутствует в обоих файлах — сравниваем рекурсивно.
			item.Status = ChangeNested
			nestedResult, err := DiffFile(map1, map2)
			if err != nil {
				return nil, err
			}
			item.Nested = nestedResult
		case isMap1 && !isMap2:
			// Мапа превратилась в скаляр: рисуем два элемента —
			// удалённую мапу и добавленный скаляр.
			removed := DiffItem{Key: k, Status: ChangeRemoved}
			nestedResult, err := diffMapAgainstEmpty(map1, true)
			if err != nil {
				return nil, err
			}
			removed.Nested = nestedResult
			items = append(items, removed)

			item.Status = ChangeAdded
			item.After = v2
		case !isMap1 && isMap2:
			// Скаляр превратился в мапу: удалённый скаляр + добавленная мапа.
			removed := DiffItem{Key: k, Status: ChangeRemoved, Before: v1}
			items = append(items, removed)

			item.Status = ChangeAdded
			nestedResult, err := diffMapAgainstEmpty(map2, false)
			if err != nil {
				return nil, err
			}
			item.Nested = nestedResult
		default:
			status, a, b := compareScalar(v1, v2)
			item.Status = status
			item.Before = a
			item.After = b
		}
		items = append(items, item)
	}
	return &DiffResult{Items: items}, nil
}

// keysCollect собирает все уникальные ключи из двух структур
func keysCollect(config1, config2 map[string]interface{}) []string {
	keysMap := make(map[string]struct{})
	for k := range config1 {
		keysMap[k] = struct{}{}
	}
	for k := range config2 {
		keysMap[k] = struct{}{}
	}

	// Превращаем ключи в срез и сортируем по алфавиту (для детерминированного вывода)
	keys := make([]string, 0, len(keysMap))
	for k := range keysMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// diffMapAgainstEmpty делает сравнение существующей мапы с пустой
func diffMapAgainstEmpty(m map[string]interface{}, isLeft bool) (*DiffResult, error) {
	empty := make(map[string]interface{})
	if isLeft {
		return DiffFile(m, empty)
	}
	return DiffFile(empty, m)
}

// compareScalar сравнивает два скаляра
func compareScalar(v1, v2 interface{}) (status string, before, after interface{}) {
	if reflect.DeepEqual(v1, v2) {
		return ChangeUnchanged, v1, nil
	}
	return ChangeUpdated, v1, v2
}
