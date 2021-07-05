package server

import (
	"reflect"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/pkg"
	"github.com/dollarkillerx/light/utils"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
var typeOfContext = reflect.TypeOf((*light.Context)(nil)).Elem()

type methodType struct {
	method       reflect.Method
	RequestType  reflect.Type
	ResponseType reflect.Type
}

// constructionMethods Get specific method
func constructionMethods(typ reflect.Type) (map[string]*methodType, error) {
	methods := make(map[string]*methodType)
	for idx := 0; idx < typ.NumMethod(); idx++ {
		method := typ.Method(idx)
		mType := method.Type
		mName := method.Name

		if !utils.IsPublic(mName) {
			return nil, pkg.ErrNonPublic
		}

		// 默认是4个
		if mType.NumIn() != 4 { // func(*server.MethodTest, *light.Context, *server.MethodTestReq, *server.MethodTestResp) error
			continue
		}

		// 检验它第一个参数是否是ctx
		ctxType := mType.In(1)
		if !(ctxType.Elem() == typeOfContext) {
			continue
		}

		// request 参数检查
		requestType := mType.In(2)
		if requestType.Kind() != reflect.Ptr {
			continue
		}

		if !utils.IsPublicOrBuiltinType(requestType) {
			continue
		}

		// response 参数检查
		responseType := mType.In(3)
		if responseType.Kind() != reflect.Ptr {
			continue
		}

		if !utils.IsPublicOrBuiltinType(responseType) {
			continue
		}

		// 校验返回参数
		if mType.NumOut() != 1 {
			continue
		}

		returnType := mType.Out(0)
		if returnType != typeOfError {
			continue
		}

		methods[mName] = &methodType{
			method:       method,
			RequestType:  requestType,
			ResponseType: responseType,
		}
	}

	if len(methods) == 0 {
		return nil, pkg.ErrNoAvailable
	}

	return methods, nil
}
