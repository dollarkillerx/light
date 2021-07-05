package main

import (
	"fmt"

	"github.com/dollarkillerx/light"
)

func main() {

}

type server struct {
}

type HelloWorldRequest struct {
	Name string
}

type HelloWorldResponse struct {
	Msg string
}

func (s *server) HelloWorld(ctx *light.Context, req *HelloWorldRequest, resp *HelloWorldResponse) error {
	resp.Msg = fmt.Sprintf("hello world by: %s", req.Name)
	return nil
}
