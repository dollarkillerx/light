package transport

import (
	"errors"
	"fmt"
	"net"
)

type transport struct {
	trMap map[string]genTransport
}

type genTransport func(addr string) (net.Listener, error)

var Transport = &transport{
	trMap: map[string]genTransport{},
}

func (t *transport) register(trType string, genTr genTransport) {
	t.trMap[trType] = genTr
}

func (t *transport) Gen(trType string, addr string) (net.Listener, error) {
	gFn, ex := t.trMap[trType]
	if !ex {
		return nil, errors.New(fmt.Sprintf("%s not funod", trType))
	}

	return gFn(addr)
}
