package server

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/dollarkillerx/light"
)

type TestMethod struct {
}

type MethodTestReq struct {
	Name string
}

type MethodTestResp struct {
	RPName string
}

func (m *TestMethod) HelloWorld(ctx *light.Context, req *MethodTestReq, resp *MethodTestResp) error {
	resp.RPName = fmt.Sprintf("hello: %s", req.Name)
	//return errors.New("what ?")
	return nil
}

func TestServer(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := NewServer()
	err := server.Register(&TestMethod{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = server.Run(Trace())
	if err != nil {
		log.Fatalln(err)
	}
}

func TestMethodF(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	mt := &TestMethod{}

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

func TestManager(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	mt := &TestMethod{}

	server := NewServer()
	err := server.Register(mt)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
