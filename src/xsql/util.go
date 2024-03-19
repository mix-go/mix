package xsql

import (
	"fmt"
	"reflect"
)

func BuildTagValues(tagKey string, ptr interface{}, pairs ...interface{}) (map[string]interface{}, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("xsql: arguments should be in pairs")
	}

	result := make(map[string]interface{})
	structValue := reflect.ValueOf(ptr).Elem()
	structType := structValue.Type()

	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xsql: ptr must be a pointer to a struct")
	}

	for i := 0; i < len(pairs); i += 2 {
		fieldPtrValue := reflect.ValueOf(pairs[i])
		if fieldPtrValue.Kind() != reflect.Ptr || fieldPtrValue.Elem().Kind() == reflect.Ptr {
			return nil, fmt.Errorf("xsql: argument at index %d must be a non-nil pointer to a struct field", i)
		}

		var fieldName string
		found := false
		for j := 0; j < structValue.NumField(); j++ {
			if structValue.Field(j).Addr().Interface() == pairs[i] {
				fieldName = structType.Field(j).Tag.Get(tagKey)
				if fieldName == "" {
					return nil, fmt.Errorf("xsql: no matching struct field tag found for pointer at index %d", i)
				}
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("xsql: no matching struct field found for pointer at index %d", i)
		}

		// Set the field name and value in the map
		result[fieldName] = pairs[i+1]
	}

	return result, nil
}
