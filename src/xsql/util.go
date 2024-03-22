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
	value := reflect.ValueOf(ptr).Elem()

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xsql: ptr must be a pointer to a struct")
	}

	fieldsMap := map[string]reflect.Value{}
	populateFieldsMap(tagKey, value, fieldsMap)

	for i := 0; i < len(pairs); i += 2 {
		fieldPtr, ok := pairs[i].(interface{})
		if !ok {
			return nil, fmt.Errorf("xsql: argument at index %d is not a pointer", i)
		}

		fieldValue := reflect.ValueOf(fieldPtr)
		if fieldValue.Kind() != reflect.Ptr || fieldValue.IsNil() {
			return nil, fmt.Errorf("xsql: argument at index %d must be a non-nil pointer to a struct field", i)
		}

		var fieldName string
		var found bool
		for name, field := range fieldsMap {
			if field.Addr().Interface() == fieldPtr {
				fieldName = name
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("xsql: no matching struct field found for pointer at index %d", i)
		}

		result[fieldName] = pairs[i+1]
	}

	return result, nil
}

// populateFieldsMap is a recursive function that maps field names to their values,
// including fields from embedded structs.
func populateFieldsMap(tagKey string, v reflect.Value, fieldsMap map[string]reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)
		tag := fieldType.Tag.Get(tagKey)
		// If it's an embedded struct, we need to recurse into it
		if fieldType.Anonymous && fieldValue.Type().Kind() == reflect.Struct {
			populateFieldsMap(tagKey, fieldValue, fieldsMap)
		} else if tag != "" {
			// Only add the field if it has the xsql tag
			fieldName := tag
			fieldsMap[fieldName] = fieldValue
		}
	}
}
