package common

import (
	"strings"
	"reflect"
)

// GetTypeName identify type name and return reflect type and name's type
func GetTypeName(something interface{}) (reflect.Type, string) {
	rawType := reflect.TypeOf(something)

	// Return value if that var is pointer
	if (rawType.Kind() == reflect.Ptr) {
		rawType = rawType.Elem()
	}

	rawTypeStr := rawType.String()
	rawTypeStrParts := strings.Split(rawTypeStr, ".") // bo phan package name
	return rawType, rawTypeStrParts[1]
}
