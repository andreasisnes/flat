package goflat

import (
	"reflect"
	"strings"
)

func Copy(data interface{}, value interface{}) {
	if data == nil || value == nil {
		return
	}

	copy(data, value)
}

func copy(data interface{}, result interface{}) interface{} {
	rdType, rdValue := reflectValue(data)
	if rdType.Kind() == reflect.Ptr {
		return copy(rdValue.Elem().Interface(), result)
	}

	switch rdType.Kind() {
	case reflect.Slice:
		for i := 0; i < rdValue.Len(); i++ {
			copy(rdValue.Index(i).Interface(), result)
		}

		return result
	case reflect.Struct:
		for i := 0; i < rdValue.NumField(); i++ {
			dName := reflect.Indirect(rdValue).Type().Field(i).Name
			copyValueFromStruct(dName, rdValue.Field(i).Interface(), result)
		}
		return result

	case reflect.Map:
		for _, rKey := range rdValue.MapKeys() {
			copy(rdValue.MapIndex(rKey).Interface(), result)
		}
		return result
	default:
		copyTo(result, data)
	}

	return result
}

func copyValueFromStruct(key string, value interface{}, result interface{}) {
	rrType, rrValue := reflectValue(result)

	if rrType.Elem().Kind() == reflect.Ptr {
		copyValueFromStruct(key, value, rrValue.Elem().Interface())
	}

	switch rrType.Elem().Kind() {
	case reflect.Struct:
		for i := 0; i < rrValue.Elem().NumField(); i++ {
			fieldName := reflect.Indirect(rrValue.Elem()).Type().Field(i).Name
			if strings.EqualFold(key, fieldName) {
				copy(value, rrValue.Elem().Field(i).Addr().Interface())
			}
		}
	}
}

// Copies a value the parameter from to the value to using reflection.
func copyTo(to interface{}, from interface{}) {
	f := reflect.ValueOf(from)
	t := reflect.ValueOf(to)

	if f.Kind() == reflect.Ptr {
		copyTo(to, f.Elem().Interface())
	}

	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()

	if t.Kind() == reflect.Ptr {
		copyTo(t.Elem().Interface(), from)
	}

	if t.IsValid() && t.CanSet() {
		t.Set(f)
	}
}

func reflectValue(value interface{}) (reflect.Type, reflect.Value) {
	return reflect.TypeOf(value), reflect.ValueOf(value)
}
