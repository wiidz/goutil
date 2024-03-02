package typeHelper

import (
	"errors"
	"reflect"
)

// StructToMap 结构体转换为map
func StructToMap(obj interface{}) (result map[string]interface{}, err error) {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	if objType.Kind() != reflect.Struct {
		err = errors.New("input object is not a struct")
		return
	}

	result = make(map[string]interface{})

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i).Interface()
		result[field.Name] = fieldValue
	}

	return
}
