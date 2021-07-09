package client

import (
	"fmt"
	"log"
	"testing"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/transport"
)

type MethodTestReq struct {
	Name string
}

type MethodTestResp struct {
	RPName string
}

func TestClient(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	client := NewClient(discovery.NewSimplePeerToPeer("127.0.0.1:8397", transport.TCP))
	connect, err := client.NewConnect("TestMethod")
	if err != nil {
		log.Fatalln(err)
		return
	}

	req := MethodTestReq{Name: "hello world"}
	resp := MethodTestResp{}
	ctx := light.DefaultCtx()
	ctx.SetValue("AUTH", "XPO")
	err = connect.Call(ctx, "HelloWorld", &req, &resp)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println("resp: ", resp)
}
