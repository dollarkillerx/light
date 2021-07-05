package server

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/dollarkillerx/light"
)

type MethodTest struct {
}

type MethodTestReq struct {
}

type MethodTestResp struct {
}

func (m *MethodTest) HelloWorld(ctx *light.Context, req *MethodTestReq, resp *MethodTestResp) error {
	return nil
}

func TestMethod(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	mt := &MethodTest{}

	mtType := reflect.TypeOf(mt)
	mtVal := reflect.ValueOf(mt)

	fmt.Println(reflect.Indirect(mtVal).Type().Name())
	methods, err := constructionMethods(mtType)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(methods)

	//l1 := &light.Context{}
	//l2 := &light.Context{}
	//fmt.Println(reflect.TypeOf(l1).Elem() == reflect.TypeOf(l2).Elem())
}
