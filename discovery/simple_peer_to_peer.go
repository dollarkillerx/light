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

func (s *SimplePeerToPeer) Registry(serName, addr string, weights float64, protocol transport.Protocol, maximumLoad int64, serID *string) error {
	s.ser = &Server{
		ServerName:  serName,
		Addr:        addr,
		ID:          *serID,
		Weights:     weights,
		Protocol:    protocol,
		MaximumLoad: maximumLoad,
	}
	return nil
}

func (s *SimplePeerToPeer) UnRegistry(serName string, serID string) error {
	return nil
}

func (s *SimplePeerToPeer) Add(load int64) {
	return
}

func (s *SimplePeerToPeer) Less(load int64) {
	return
}

func (s *SimplePeerToPeer) Limit() bool {
	return false
}
