package server

import "reflect"

type service struct {
	name       string                 // server name
	refVal     reflect.Value          // server reflect value
	refType    reflect.Type           // server reflect type
	methodType map[string]*methodType // server method
}
