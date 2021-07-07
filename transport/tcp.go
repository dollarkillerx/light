package transport

import "net"

func init() {
	Transport.register("tcp", defaultTcp)
}

func defaultTcp(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
