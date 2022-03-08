package goflat

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrEmptyDelimiter = errors.New("delimiter can't be an empty string")
)

func UList(flatmap map[string]interface{}, options *Options) ([]interface{}, error) {
	wrapper := make(map[string]interface{})
	for key, value := range flatmap {
		wrapper[concatKey("wrapper", key, &Options{
			Delimiter: ".",
		})] = value
	}

	if result, err := UMap(wrapper, options); err == nil {
		return result["wrapper"].([]interface{}), nil
	} else {
		return nil, err
	}
}

func UMap(flatmap map[string]interface{}, options *Options) (result map[string]interface{}, err error) {
	value, err := unflatWrapper(flatmap, make(map[string]interface{}), options)
	if result, ok := value.(map[string]interface{}); ok {
		return result, err
	}

	return nil, err
}

func unflatWrapper(flat map[string]interface{}, result interface{}, options *Options) (interface{}, error) {
	if options == nil {
		options = &Options{
			Delimiter: DefaultDelimiter,
		}
	}

	if options.Delimiter == "" {
		return result, ErrEmptyDelimiter
	}

	keyArrays := sortKeys(flat, options)
	for keyArrayIdx, keyArray := range keyArrays {
		unflat(keyArrays, result, flat[strings.Join(keyArray, options.Delimiter)], keyArray, keyArrayIdx, 0)
	}

	return result, nil
}

func unflat(keyArrays [][]string, currValue interface{}, value interface{}, path []string, keyArrayIdx, depth int) interface{} {
	switch m := currValue.(type) {
	case []interface{}:
		return unflatList(keyArrays, m, value, path, keyArrayIdx, depth)
	case map[string]interface{}:
		return unflatObject(keyArrays, m, value, path, keyArrayIdx, depth)
	default:
		panic(fmt.Sprintf("can't parse child of type: %s.", reflect.TypeOf(m).Name()))
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
	if depth >= len(currKey) {
		return false
	}
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
