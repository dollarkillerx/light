package serialization_plugin

import (
	"fmt"
	"reflect"
)

type ByteCode struct{}

func (j *ByteCode) Encode(i interface{}) ([]byte, error) {
	if data, ok := i.([]byte); ok {
		return data, nil
	}
	if data, ok := i.(*[]byte); ok {
		return *data, nil
	}

	return nil, fmt.Errorf("%T is not a []byte", i)
}

func (j *ByteCode) Decode(data []byte, i interface{}) error {
	reflect.Indirect(reflect.ValueOf(i)).SetBytes(data)
	return nil
}
