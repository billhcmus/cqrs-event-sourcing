package utils

import (
	"strings"
	"reflect"
)

// GetTypeName identify type name and return reflect type and name's type
func GetTypeName(myvar interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(myvar)

	// Neu myvar co kieu con tro thi chuyen qua value
	if (rawType.Kind() == reflect.Ptr) {
		rawType = rawType.Elem()
	}

	rawTypeStr := rawType.String()
	rawTypeStrParts := strings.Split(rawTypeStr, ".") // bo phan package name
	return rawType, rawTypeStrParts[1]
}
