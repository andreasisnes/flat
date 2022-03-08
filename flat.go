package goflat

import (
	"fmt"
)

func List(value []interface{}, options *Options) (map[string]interface{}, error) {
	return flatwrapper(value, options)
}

func Map(value map[string]interface{}, options *Options) (map[string]interface{}, error) {
	return flatwrapper(value, options)
}

func flatwrapper(unflat interface{}, options *Options) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if options == nil {
		options = &Options{
			Delimiter: DefaultDelimiter,
		}
	}

	if options.Delimiter == "" {
		return result, ErrEmptyDelimiter
	}

	return flat(unflat, result, "", options)
}

func flat(curr interface{}, result map[string]interface{}, key string, options *Options) (map[string]interface{}, error) {
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

	return result, nil
}

func concatKey(key, idx string, options *Options) string {
	if key == "" {
		return idx
	}

	return fmt.Sprintf("%s%s%s", key, options.Delimiter, idx)
}
