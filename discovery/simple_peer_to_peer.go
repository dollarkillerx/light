package discovery

import (
	"github.com/dollarkillerx/light/transport"
)

type SimplePeerToPeer struct {
	ser *Server
}

func (s *SimplePeerToPeer) Discovery(serName string) ([]*Server, error) {
	var sr []*Server
	sr = append(sr, s.ser)
	return sr, nil
}

func NewSimplePeerToPeer(addr string, protocol transport.Protocol) *SimplePeerToPeer {
	return &SimplePeerToPeer{
		ser: &Server{
			Addr:     addr,
			Protocol: protocol,
		},
	}
}

func (s *SimplePeerToPeer) Registry(serName, addr string, weights float64, protocol transport.Protocol, serID *string) error {
	return nil
}

func (s *SimplePeerToPeer) UnRegistry(serName string, serID string) error {
	return nil
}
