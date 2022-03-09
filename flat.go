package goflat

import (
	"fmt"
	"reflect"
)

// This function takes a map and flattens it by assigning each leaf value to a unique
// concatenated key that represents the value's path. A delimiter always gets in between
// keys that get concatenated. If the options parameter is nil or delimiter is an empty
// string this method will use the default delimiter ".".
func Map(value interface{}, options *Options) map[string]interface{} {
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
	rType := reflect.TypeOf(curr)
	rValue := reflect.ValueOf(curr)

	if rType.Kind() == reflect.Ptr {
		flat(rValue.Elem().Interface(), result, key, options)
	}

	switch rType.Kind() {
	case reflect.Slice:
		for i := 0; i < rValue.Len(); i++ {
			flat(rValue.Index(i).Interface(), result, concatKey(key, fmt.Sprintf("%d", i), options), options)
		}
	case reflect.Struct:
		for i := 0; i < rValue.NumField(); i++ {
			name := reflect.Indirect(rValue).Type().Field(i).Name
			flat(rValue.Field(i).Interface(), result, concatKey(key, name, options), options)
		}
	case reflect.Map:
		for _, rKey := range rValue.MapKeys() {
			flat(rValue.MapIndex(rKey).Interface(), result, concatKey(key, fmt.Sprintf("%v", rKey.Interface()), options), options)
		}
	default:
		if key != "" {
			result[getFold(key, options)] = curr
		}
	}

	return result
}
