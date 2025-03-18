package golangutils

import (
	"reflect"
	"strings"
	"time"
)

func Mapping[T interface{}](source interface{}) T {
	var target T
	MappingByReflect(reflect.ValueOf(source), reflect.ValueOf(&target))

	return target
}

func MappingByReflect(source reflect.Value, target reflect.Value) {
	source = IndirectValue(source)
	target = IndirectValue(target)
	if source.Kind() == reflect.Struct && target.Kind() == reflect.Struct {
		targetFieldsMap := make(map[string]reflect.Value)

		for i := 0; i < target.NumField(); i++ {
			targetField := target.Field(i)
			targetFieldsMap[target.Type().Field(i).Name] = targetField
		}

		for i := 0; i < target.NumField(); i++ {
			aliasString := target.Type().Field(i).Tag.Get("mappingalias")
			if aliasString != "" {
				aliases := strings.Split(aliasString, ",")
				targetField := target.Field(i)
				for _, alias := range aliases {
					if _, ok := targetFieldsMap[alias]; !ok {
						targetFieldsMap[alias] = targetField
					}
				}
			}
		}

		for i := 0; i < source.NumField(); i++ {
			sourceField := source.Field(i)

			if targetField, ok := targetFieldsMap[source.Type().Field(i).Name]; ok {
				if (sourceField.Kind() == reflect.Struct || sourceField.Kind() == reflect.Ptr) && (targetField.Kind() == reflect.Struct || targetField.Kind() == reflect.Ptr) {
					if targetField.Type() == reflect.TypeOf(time.Time{}) {
						if targetField.Type().ConvertibleTo(sourceField.Type()) {
							targetField.Set(sourceField.Convert(targetField.Type()))
						}
					} else {
						MappingByReflect(sourceField, targetField)
					}
				} else if sourceField.Kind() == targetField.Kind() {
					if sourceField.Type().PkgPath() == targetField.Type().PkgPath() && sourceField.Type().Name() == targetField.Type().Name() {
						targetField.Set(sourceField)
					} else {
						if targetField.Type().ConvertibleTo(sourceField.Type()) {
							targetField.Set(sourceField.Convert(targetField.Type()))
						}
					}
				}
			}
		}
	}
}
