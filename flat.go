package goflat

import (
	"fmt"
)

// This function takes a list and flattens it by assigning each leaf value to a unique
// concatenated key that represents the value's path. A delimiter always gets in between
// keys that get concatenated. If the options parameter is nil or delimiter is an empty
// string this method will use the default delimiter ".".
func List(value []interface{}, options *Options) map[string]interface{} {
	if value == nil {
		return nil
	}

	return flatwrapper(value, options)
}

// This function takes a map and flattens it by assigning each leaf value to a unique
// concatenated key that represents the value's path. A delimiter always gets in between
// keys that get concatenated. If the options parameter is nil or delimiter is an empty
// string this method will use the default delimiter ".".
func Map(value map[string]interface{}, options *Options) map[string]interface{} {
	if value == nil {
		return nil
	}

	return flatwrapper(value, options)
}

func flatwrapper(unflat interface{}, options *Options) map[string]interface{} {
	result := make(map[string]interface{})
	options = createDefultOptionsIfNil(options)

	return flat(unflat, result, "", options)
}

func flat(curr interface{}, result map[string]interface{}, key string, options *Options) map[string]interface{} {
	switch m := curr.(type) {
	case []interface{}:
		for idx, value := range m {
			flat(value, result, concatKey(key, fmt.Sprintf("%d", idx), options), options)
		}
	case map[string]interface{}:
		for idx, value := range m {
			flat(value, result, concatKey(key, idx, options), options)
		}
	default:
		result[key] = curr
	}

	return result
}

func concatKey(key, idx string, options *Options) string {
	if key == "" {
		return idx
	}

	return fmt.Sprintf("%s%s%s", key, options.Delimiter, idx)
}
