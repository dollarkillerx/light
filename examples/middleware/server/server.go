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

	ser.Before(func(ctx *light.Context, request interface{}, response interface{}) error {
		fmt.Println("before 1")
		ctx.SetValue("s1", "s1")
		return nil
	})

	ser.After(func(ctx *light.Context, request interface{}, response interface{}) error {
		fmt.Println("after 1")
		ctx.SetValue("s5", "s5")
		return nil
	})

	ser.BeforePath("helloworld.HelloWorld", func(ctx *light.Context, request interface{}, response interface{}) error {
		fmt.Println("before 2")
		ctx.SetValue("s2", "s2")
		return nil
	})

	ser.AfterPath("helloworld.HelloWorld", func(ctx *light.Context, request interface{}, response interface{}) error {
		fmt.Println("after 2")
		fmt.Println(ctx.GetMetaData())
		return nil
	})

	if err := ser.Run(server.UseTCP("0.0.0.0:8074"), server.Trace()); err != nil {
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
	ctx.SetValue("s3", "s3")
	resp.RPName = fmt.Sprintf("hello world by: %s", req.Name)
	//return errors.New(":xx")
	//fmt.Println(resp)
	return nil
}
