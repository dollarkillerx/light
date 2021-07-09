package transport

import "net"

func init() {
	Transport.register(TCP, defaultTcp)
	Client.register(TCP, defaultTcpClient)
}

func defaultTcp(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}

func defaultTcpClient(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}
