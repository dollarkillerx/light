package transport

import (
	"net"

	"github.com/xtaci/kcp-go/v5"
)

func init() {
	Transport.register(KCP, kcpNet)
	Client.register(KCP, kcpNetClient)
}

func kcpNet(addr string) (net.Listener, error) {
	return kcp.Listen(addr)
}

func kcpNetClient(addr string) (net.Conn, error) {
	return kcp.Dial(addr)
}
