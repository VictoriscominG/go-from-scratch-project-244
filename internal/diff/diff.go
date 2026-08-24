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

	// Собираем все уникальные ключи из двух структур
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
			if isMap2 {
				// Добавили мапу целиком: внутри неё ничего не сравнивается,
				// все вложенные элементы помечаем как added (со значениями).
				emptyMap := make(map[string]interface{})
				nestedResult, err := DiffFile(emptyMap, map2)
				if err != nil {
					return nil, err
				}
				item.Status = ChangeAdded
				item.Nested = nestedResult
			} else {
				// Добавленный скаляр
				item.Status = ChangeAdded
				item.After = v2
			}
		case ok1 && !ok2:
			// Удалено
			if isMap1 {
				// Удалили мапу целиком: все вложенные элементы — removed (со значениями).
				emptyMap := make(map[string]interface{})
				nestedResult, err := DiffFile(map1, emptyMap)
				if err != nil {
					return nil, err
				}
				item.Status = ChangeRemoved
				item.Nested = nestedResult
			} else {
				// Удалённый скаляр
				item.Status = ChangeRemoved
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
			emptyMap := make(map[string]interface{})
			nestedResult, err := DiffFile(map1, emptyMap)
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
			emptyMap := make(map[string]interface{})
			nestedResult, err := DiffFile(emptyMap, map2)
			if err != nil {
				return nil, err
			}
			item.Nested = nestedResult
		default:
			// Сравнение двух скаляров
			if reflect.DeepEqual(v1, v2) {
				// Значение не изменилось
				item.Status = ChangeUnchanged
				item.Before = v1
			} else {
				// Значение изменилось
				item.Status = ChangeUpdated
				item.Before = v1
				item.After = v2
			}
		}
		items = append(items, item)
	}
	return &DiffResult{Items: items}, nil
}
