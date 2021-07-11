package main

import (
	"fmt"
	"log"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/server"
)

func main() {
	ser := server.NewServer()
	err := ser.RegisterName(&HelloWorld{}, "helloworld")
	if err != nil {
		log.Fatalln(err)
	}

	if err := ser.Run(server.UseTCP("0.0.0.0:8074"), server.Trace(), server.SetAESCryptology([]byte("58a95a8f804b49e686f651a0d3f6e631"))); err != nil {
		log.Fatalln(err)
	}
}

type HelloWorld struct{}

type HelloWorldRequest struct {
	Name string
}

type HelloWorldResponse struct {
	RPName string
}

func (s *HelloWorld) HelloWorld(ctx *light.Context, req *HelloWorldRequest, resp *HelloWorldResponse) error {
	resp.RPName = fmt.Sprintf("hello world by: %s", req.Name)
	//return errors.New(":xx")
	//fmt.Println(resp)
	return nil
}
