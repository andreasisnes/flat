package flat

import (
	"fmt"
	"reflect"
)

// Map takes a value of any type and flattens it by assigning each leaf value to a unique
// concatenated key that represents the value's path. A delimiter is used to separate keys
// that are concatenated. If the options parameter is nil or delimiter is an empty string,
// the method will use the default delimiter ".".
//
// Parameters:
//   - input: The input value to be flattened.
//   - opts: Optional configuration functions for customizing the mapping process.
//
// Returns:
//
//	A map where keys represent the paths of the original values, and values are the leaf values.
//
// Example Usage:
//
//	type Person struct {
//	    Name struct {
//	        First string
//	        Last  string
//	    }
//	    Age int
//	}
//
//	func main() {
//	    // Create a complex structure
//	    complexData := Person{
//	        Name: struct {
//	            First string
//	            Last  string
//	        }{
//	            First: "John",
//	            Last:  "Doe",
//	        },
//	        Age: 30,
//	    }
//
//	    // Flatten the complex structure
//	    flatMap := flat.Map(complexData, flat.WithDelimiter("_"))
//
//	    // Print the flattened map
//	    fmt.Println("Flattened Map:")
//	    for key, value := range flatMap {
//	        fmt.Printf("%s: %v\n", key, value)
//	    }
//	}
//
// This will flatten the complex structure `complexData` using "_" as the delimiter.
func Map(input any, opts ...func(options *Options)) map[string]any {
	if input == nil {
		return nil
	}

	options := DefaultOptions(opts...)
	return flatwrapper(input, options)
}

func flatwrapper(unflat any, options *Options) map[string]any {
	result := make(map[string]any)
	return flat(unflat, result, "", options)
}

func flat(curr any, result map[string]any, key string, options *Options) map[string]any {
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
			result[options.Fold(key)] = curr
		}
	}

	return result
}
