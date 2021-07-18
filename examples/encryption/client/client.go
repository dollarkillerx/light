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
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMacmioOUq4HTMKVutxsrWimQO
vFOZIU93NKJYusRFV8lN8NB0dSg5AcbZgyYegY07mWXBBg8zlI+4PphUj40kn0F3
aOnvZ6WrsyYlPi1ZnXBaTFXxC6YN2LH9Lb9KaWrtZH4AM+6PoXIjmtWhpZr1JuuD
7J540DLMsuoEGzLoRQIDAQAB
-----END PUBLIC KEY-----`)

	client := client.NewClient(discovery.NewSimplePeerToPeer("127.0.0.1:8074", transport.TCP), client.SetRASPublicKey(publicKey))
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
}
