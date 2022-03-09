package goflat

import (
	"sort"
	"strconv"
	"strings"
)

// UMap takes a flattened map and recreates it with nested objects and lists.
// Remember to use the same delimiter when you first flattened it.
// If the options parameter is nil or delimiter is an empty string this method will use the default delimiter ".".
func UMap(flatmap map[string]interface{}, options *Options) interface{} {
	if len(flatmap) == 0 {
		return nil
	}

	options = createDefultOptionsIfNil(options)
	keyArrays := sortKeys(flatmap, options)
	if isList(keyArrays, keyArrays[0], 0, 0) {
		wrapper := make(map[string]interface{})

		for key, value := range flatmap {
			wrapper[concatKey("wrapper", key, options)] = value
		}

		return unflatWrapper(wrapper, make(map[string]interface{}), sortKeys(wrapper, options), options).(map[string]interface{})["wrapper"]
	}

	return unflatWrapper(flatmap, make(map[string]interface{}), keyArrays, options)
}

func unflatWrapper(flat map[string]interface{}, result interface{}, keyArrays [][]string, options *Options) interface{} {
	for keyArrayIdx, keyArray := range keyArrays {
		unflat(keyArrays, result, flat[strings.Join(keyArray, options.Delimiter)], keyArray, keyArrayIdx, 0)
	}

	return result
}

func unflat(keyArrays [][]string, currValue interface{}, value interface{}, path []string, keyArrayIdx, depth int) interface{} {
	switch m := currValue.(type) {
	case []interface{}:
		return unflatList(keyArrays, m, value, path, keyArrayIdx, depth)
	default:
		return unflatObject(keyArrays, m.(map[string]interface{}), value, path, keyArrayIdx, depth)
	}
}

func unflatList(keyArrays [][]string, currValue []interface{}, value interface{}, path []string, keyArrayIdx, depth int) interface{} {
	if _, err := strconv.Atoi(path[depth]); err == nil && len(path[depth:]) == 1 {
		return append(currValue, value)
	}

	idx, _ := strconv.Atoi(path[depth])
	if idx >= 0 && idx < len(currValue) {
		currValue[idx] = unflat(keyArrays, currValue[idx], value, path, keyArrayIdx, depth+1)
		return currValue
	}

	if isList(keyArrays, path, keyArrayIdx, depth+1) {
		return append(currValue, unflat(keyArrays, make([]interface{}, 0), value, path, keyArrayIdx, depth+1))
	}

	return append(currValue, unflat(keyArrays, make(map[string]interface{}), value, path, keyArrayIdx, depth+1))
}

func unflatObject(keyArrays [][]string, currValue map[string]interface{}, value interface{}, path []string, keyArrayIdx, depth int) interface{} {
	if (depth + 1) >= len(path) {
		currValue[path[depth]] = value
		return currValue
	}

	if isList(keyArrays, path, keyArrayIdx, depth+1) {
		if val, ok := currValue[path[depth]]; ok {
			currValue[path[depth]] = unflat(keyArrays, val, value, path, keyArrayIdx, depth+1)
		} else {
			currValue[path[depth]] = unflat(keyArrays, make([]interface{}, 0), value, path, keyArrayIdx, depth+1)
		}
		return currValue
	}

	if val, ok := currValue[path[depth]]; ok {
		currValue[path[depth]] = unflat(keyArrays, val, value, path, keyArrayIdx, depth+1)
	} else {
		currValue[path[depth]] = unflat(keyArrays, make(map[string]interface{}), value, path, keyArrayIdx, depth+1)
	}

	return currValue
}

func isList(keyArrays [][]string, currKey []string, keyArrayIdx, depth int) bool {
	if _, err := strconv.Atoi(currKey[depth]); err != nil {
		return false
	}

	for idx, keyArray := range keyArrays {
		if idx == keyArrayIdx {
			continue
		}
		for i := 0; i <= depth; i++ {
			if i == depth && len(keyArray) >= depth {
				if _, err := strconv.Atoi(keyArray[i]); err != nil {
					return false
				}
			}

			if !strings.EqualFold(keyArray[i], currKey[i]) {
				break
			}
		}
	}

	return true
}

func sortKeys(flatmap map[string]interface{}, options *Options) [][]string {
	keys := make([]string, len(flatmap))
	i := 0
	for k := range flatmap {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	result := make([][]string, len(keys))
	for idx, key := range keys {
		result[idx] = strings.Split(key, options.Delimiter)
	}

	return result
}
