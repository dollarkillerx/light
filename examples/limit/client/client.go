package main

import (
	"log"
	"time"

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
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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
	//req := MethodTestReq{
	//	Name: "hello",
	//}
	//resp := MethodTestResp{}
	//err = connect.Call(light.DefaultCtx(), "HelloWorld", &req, &resp)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//for {
	//	select {
	//
	//	}
	//}

	li := make(chan struct{}, 1)
	for {
		li <- struct{}{}
		go func() {
			defer func() {
				<-li
			}()

			req := MethodTestReq{
				Name: "hello",
			}
			resp := MethodTestResp{}
			err = connect.Call(light.DefaultCtx(), "HelloWorld", &req, &resp)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(resp)
		}()

		time.Sleep(time.Millisecond * 10)
	}

}
