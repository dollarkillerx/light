package client

import (
	"net"

	"github.com/dollarkillerx/light"
)

type BaseClient struct {
	conn net.Conn
}

func newBaseClient(con net.Conn) *BaseClient {
	return &BaseClient{
		conn: con,
	}
}

func (b *BaseClient) Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) error {
	return nil
}
