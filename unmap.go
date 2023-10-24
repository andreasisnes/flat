package flat

import (
	"sort"
	"strconv"
	"strings"
)

// UMap takes a flattened map and recreates it with nested objects and lists.
// To ensure correct reconstruction, use the same delimiter as when the map was originally flattened.
// If the options parameter is nil or delimiter is an empty string,
// this method will use the default delimiter ".".
//
// Parameters:
//   - flatmap: The flattened map to be reconstructed.
//   - opts: Optional configuration functions for customizing the unflattening process.
//
// Returns:
//
//	The reconstructed nested structure from the flattened map.
//
// Example Usage:
//
//	flatMap := map[string]any{
//	  "person.name.first": "John",
//	  "person.name.last":  "Doe",
//	  "person.age":        30,
//	}
//	nestedStructure := flat.UMap(flatMap)
//
// This will reconstruct a nested structure from the flat map.
func UMap(flatmap map[string]any, opts ...func(options *Options)) any {
	if len(flatmap) == 0 {
		return nil
	}

	options := DefaultOptions(opts...)
	keyArrays := sortKeys(flatmap, options)
	if isList(keyArrays, keyArrays[0], 0, 0) {
		wrapper := make(map[string]any)
		for key, value := range flatmap {
			wrapper[concatKey("wrapper", key, options)] = value
		}

		return unflatWrapper(wrapper, make(map[string]any), sortKeys(wrapper, options), options).(map[string]any)["wrapper"]
	}

	return unflatWrapper(flatmap, make(map[string]any), keyArrays, options)
}

func unflatWrapper(flat map[string]any, result any, keyArrays [][]string, options *Options) any {
	for keyArrayIdx, keyArray := range keyArrays {
		unflat(keyArrays, result, flat[strings.Join(keyArray, options.Delimiter)], keyArray, keyArrayIdx, 0)
	}

	return result
}

func unflat(keyArrays [][]string, currValue any, value any, path []string, keyArrayIdx, depth int) any {
	switch m := currValue.(type) {
	case []any:
		return unflatList(keyArrays, m, value, path, keyArrayIdx, depth)
	default:
		return unflatObject(keyArrays, m.(map[string]any), value, path, keyArrayIdx, depth)
	}
}

func unflatList(keyArrays [][]string, currValue []any, value any, path []string, keyArrayIdx, depth int) any {
	if _, err := strconv.Atoi(path[depth]); err == nil && len(path[depth:]) == 1 {
		return append(currValue, value)
	}

	idx, _ := strconv.Atoi(path[depth])
	if idx >= 0 && idx < len(currValue) {
		currValue[idx] = unflat(keyArrays, currValue[idx], value, path, keyArrayIdx, depth+1)
		return currValue
	}

	if isList(keyArrays, path, keyArrayIdx, depth+1) {
		return append(currValue, unflat(keyArrays, make([]any, 0), value, path, keyArrayIdx, depth+1))
	}

	return append(currValue, unflat(keyArrays, make(map[string]any), value, path, keyArrayIdx, depth+1))
}

func unflatObject(keyArrays [][]string, currValue map[string]any, value any, path []string, keyArrayIdx, depth int) any {
	if (depth + 1) >= len(path) {
		currValue[path[depth]] = value
		return currValue
	}

	if isList(keyArrays, path, keyArrayIdx, depth+1) {
		if val, ok := currValue[path[depth]]; ok {
			currValue[path[depth]] = unflat(keyArrays, val, value, path, keyArrayIdx, depth+1)
		} else {
			currValue[path[depth]] = unflat(keyArrays, make([]any, 0), value, path, keyArrayIdx, depth+1)
		}
		return currValue
	}

	if val, ok := currValue[path[depth]]; ok {
		currValue[path[depth]] = unflat(keyArrays, val, value, path, keyArrayIdx, depth+1)
	} else {
		currValue[path[depth]] = unflat(keyArrays, make(map[string]any), value, path, keyArrayIdx, depth+1)
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

func sortKeys(flatmap map[string]any, options *Options) [][]string {
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
