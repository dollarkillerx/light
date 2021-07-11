package server

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/utils"
)

type service struct {
	name       string                 // server name
	refVal     reflect.Value          // server reflect value
	refType    reflect.Type           // server reflect type
	methodType map[string]*methodType // server method
}

func newService(server interface{}, serverName string, useName bool) (*service, error) {
	ser := &service{
		refVal:  reflect.ValueOf(server),
		refType: reflect.TypeOf(server),
	}

	sName := reflect.Indirect(ser.refVal).Type().Name()
	if !utils.IsPublic(sName) {
		return nil, pkg.ErrNonPublic
	}

	if useName {
		if serverName == "" {
			return nil, errors.New("Server Name is null")
		}

		sName = serverName
	}

	ser.name = sName
	methods, err := constructionMethods(ser.refType)
	if err != nil {
		return nil, err
	}
	ser.methodType = methods

	for _, v := range methods {
		log.Println("Registry Service: ", ser.name, "   method: ", v.method.Name)
	}

	return ser, nil
}

// call 方法调用
func (s *service) call(ctx *light.Context, mType *methodType, request, response reflect.Value) (err error) {
	// recover 捕获堆栈消息
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[painc service internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mType.method.Name, request.Interface(), buf)
			log.Println(err)
		}
	}()

	fn := mType.method.Func
	returnValue := fn.Call([]reflect.Value{s.refVal, reflect.ValueOf(ctx), request, response})
	errInterface := returnValue[0].Interface()
	if errInterface != nil {
		return errInterface.(error)
	}

	return nil
}
