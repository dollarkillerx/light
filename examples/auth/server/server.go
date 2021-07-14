package main

import (
	"errors"
	"fmt"
	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/server"
	"log"
)

func main() {
	ser := server.NewServer()
	err := ser.RegisterName(&HelloWorld{}, "helloworld")
	if err != nil {
		log.Fatalln(err)
	}

	if err := ser.Run(server.UseTCP("0.0.0.0:8074"), server.SetAUTH(authFunc)); err != nil {
		log.Fatalln(err)
	}
}

func authFunc(ctx *light.Context, token string) error {
	if token == "token" {
		return nil
	}

	return errors.New("error 401")
}

type HelloWorld struct{}

type HelloWorldRequest struct {
	Name string
}

type HelloWorldResponse struct {
	Msg string
}

func (s *HelloWorld) HelloWorld(ctx *light.Context, req *HelloWorldRequest, resp *HelloWorldResponse) error {
	resp.Msg = fmt.Sprintf("hello world by: %s", req.Name)
	return nil
}
