package client

import (
	"github.com/dollarkillerx/light"
)

type Payload struct {
	ctx           *light.Context
	serviceMethod string
	request       interface{}
	response      interface{}

	respChan chan error // 通知请求已经结束了  (1.单纯关闭 说明 没有问题, 2.收到error消息说明返回错误)
}

type LightClient interface {
	Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) error
	Error() error
}
