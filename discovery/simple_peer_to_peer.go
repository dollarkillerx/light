package discovery

type SimplePeerToPeer struct {
	addr string
}

func (s SimplePeerToPeer) Discovery(serName string) ([]string, error) {
	return []string{s.addr}, nil
}

func NewSimplePeerToPeer(addr string) *SimplePeerToPeer {
	return &SimplePeerToPeer{
		addr: addr,
	}
}
