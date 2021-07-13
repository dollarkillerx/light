package transport

import "net"

func init() {
	Transport.register(UNIX, unixNet)
	Client.register(UNIX, unixNetClient)
}

func unixNet(addr string) (net.Listener, error) {
	return net.Listen("unix", addr)
}

func unixNetClient(addr string) (net.Conn, error) {
	return net.Dial("unix", addr)
}
