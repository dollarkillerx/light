package transport

import (
	"net"

	"github.com/xtaci/kcp-go/v5"
)

func init() {
	Transport.register("kcp", kcpNet)
}

func kcpNet(addr string) (net.Listener, error) {
	return kcp.Listen(addr)
}
