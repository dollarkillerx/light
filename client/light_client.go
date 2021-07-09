package client

import (
	"github.com/dollarkillerx/light"
)

type LightClient interface {
	Call(ctx *light.Context, serviceMethod string, request interface{}, response interface{}) error
}
