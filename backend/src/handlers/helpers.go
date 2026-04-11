package handlers

import (
	"reflect"
)

func getJSONFieldName(structType reflect.Type, fieldName string) string {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Name == fieldName {
			tag := field.Tag.Get("json")
			if tag != "" && tag != "-" {
				parts := splitTag(tag)
				return parts[0]
			}
			return fieldName
		}
	}
	return fieldName
}

func splitTag(tag string) []string {
	var parts []string
	current := ""
	for _, c := range tag {
		if c == ',' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}
