package utils

import (
	"fmt"
	"reflect"
	"runtime"
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

// RefNew 通过refType 构造实例
func RefNew(refType reflect.Type) interface{} {
	var refValue reflect.Value
	if refType.Kind() == reflect.Ptr {
		refValue = reflect.New(refType.Elem())
	} else {
		refValue = reflect.New(refType)
	}

	return refValue.Interface()
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}
