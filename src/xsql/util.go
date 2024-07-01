package xsql

import (
	"fmt"
	"reflect"
)

type TagValues []TagValue

type TagValue struct {
	Key   interface{}
	Value interface{}
}

// TagValuesMap takes a tag key, a pointer to a struct, and TagValues.
// It constructs a map where each key is the struct field's tag value, paired with the corresponding value from TagValues.
func TagValuesMap(tagKey string, ptr interface{}, values TagValues) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	structValue := reflect.ValueOf(ptr).Elem()

	if structValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xsql: ptr must be a pointer to a struct")
	}

	fieldsMap := make(map[string]reflect.Value)
	populateFieldsMap(tagKey, structValue, fieldsMap)

	for i, tagValue := range values {
		fieldPtr, fieldValue := tagValue.Key, reflect.ValueOf(tagValue.Key)
		if fieldValue.Kind() != reflect.Ptr || fieldValue.IsNil() {
			return nil, fmt.Errorf("xsql: error at item %d in values slice: key is not a non-nil pointer to a struct field", i)
		}

		foundFieldName := ""
		for tagName, field := range fieldsMap {
			if field.Addr().Interface() == fieldPtr {
				foundFieldName = tagName
				break
			}
		}

		if foundFieldName == "" {
			return nil, fmt.Errorf("xsql: no matching struct field found for item %d in values slice", i)
		}

		result[foundFieldName] = tagValue.Value
	}

	return result, nil
}

func populateFieldsMap(tagKey string, v reflect.Value, fieldsMap map[string]reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		tag := fieldType.Tag.Get(tagKey)
		if fieldType.Anonymous && field.Type().Kind() == reflect.Struct {
			populateFieldsMap(tagKey, field, fieldsMap)
		} else if tag != "" {
			fieldsMap[tag] = field
		}
	}
}
