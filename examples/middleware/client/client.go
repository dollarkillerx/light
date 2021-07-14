package main

import (
	"fmt"
	"log"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/client"
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/transport"
)

type MethodTestReq struct {
	Name string
}

type MethodTestResp struct {
	RPName string
}

func main() {
	client := client.NewClient(discovery.NewSimplePeerToPeer("127.0.0.1:8074", transport.TCP))
	connect, err := client.NewConnect("helloworld")
	if err != nil {
		log.Fatalln(err)
		return
	}

	req := MethodTestReq{
		Name: "hello",
	}
	resp := MethodTestResp{}
	err = connect.Call(light.DefaultCtx(), "HelloWorld", &req, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)

	for {
		select {}
	}
}
