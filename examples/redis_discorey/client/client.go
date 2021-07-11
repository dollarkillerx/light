package main

import (
	"fmt"
	"log"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/client"
	"github.com/dollarkillerx/light/discovery"
)

type MethodTestReq struct {
	Name string
}

type MethodTestResp struct {
	RPName string
}

func main() {
	redisDiscovery, err := discovery.NewRedisDiscovery("127.0.0.1:6379", 10, nil)
	if err != nil {
		log.Fatalln(err)
	}

	c := client.NewClient(redisDiscovery)
	connect, err := c.NewConnect("helloworld")
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
