package discovery

type SimplePeerToPeer struct {
	addr string
}

func (s SimplePeerToPeer) Discovery(serName string) []string {
	return []string{s.addr}
}

func NewSimplePeerToPeer(addr string) *SimplePeerToPeer {
	return &SimplePeerToPeer{
		addr: addr,
	}
}
