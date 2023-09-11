package xconv

import (
	"reflect"
	"strings"
	"unsafe"
)

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StructToMap Convert struct to map.
// This function first tries to use the bson tag, and if the bson tag does not exist, it will use the json tag.
// if both bson and json tags do not exist, then it will use the field name as the key. Additionally,
// if the tag value is "-", this field will be ignored and not added to the map.
func StructToMap(i interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	inputType := reflect.TypeOf(i)
	inputVal := reflect.ValueOf(i)

	if inputType.Kind() == reflect.Ptr {
		inputType = inputType.Elem()
		inputVal = inputVal.Elem()
	}

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		value := inputVal.Field(i).Interface()

		key := field.Tag.Get("bson")
		if key == "" {
			key = field.Tag.Get("json")
		}

		if key == "-" {
			continue
		}

		key = strings.Split(key, ",")[0]
		if key == "" {
			key = field.Name
		}

		result[key] = value
	}

	return result
}
