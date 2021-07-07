package transport

import "net"

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
