package utils

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// IsPublic is public
func IsPublic(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func IsPublicOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return IsPublic(t.Name()) || t.PkgPath() == ""
}
