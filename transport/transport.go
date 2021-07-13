package transport

import (
	"errors"
	"fmt"
	"net"
)

type transport struct {
	trMap map[Protocol]genTransport
}

type genTransport func(addr string) (net.Listener, error)

var Transport = &transport{
	trMap: map[Protocol]genTransport{},
}

func (t *transport) register(trType Protocol, genTr genTransport) {
	t.trMap[trType] = genTr
}

func (t *transport) Gen(trType Protocol, addr string) (net.Listener, error) {
	gFn, ex := t.trMap[trType]
	if !ex {
		return nil, errors.New(fmt.Sprintf("%s not funod", trType))
	}

	return gFn(addr)
}

type transportClient struct {
	trMap map[Protocol]genTransportClient
}

type genTransportClient func(addr string) (net.Conn, error)

var Client = &transportClient{
	trMap: map[Protocol]genTransportClient{},
}

func (t *transportClient) register(trType Protocol, genTr genTransportClient) {
	t.trMap[trType] = genTr
}

func (t *transportClient) Gen(trType Protocol, addr string) (net.Conn, error) {
	gFn, ex := t.trMap[trType]
	if !ex {
		return nil, errors.New(fmt.Sprintf("%s not funod", trType))
	}

	return gFn(addr)
}

type Protocol string

const (
	TCP  Protocol = "tcp"
	KCP  Protocol = "kcp"
	MQTT Protocol = "mqtt"
	UNIX Protocol = "unix"
)

func (p Protocol) String() string {
	return string(p)
}
