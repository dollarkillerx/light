package discovery

import (
	"github.com/dollarkillerx/light/transport"
)

// Discovery 服务发现
type Discovery interface {
	Discovery(serName string) ([]*Server, error)
	Registry(serName, addr string, weights float64, protocol transport.Protocol, serID *string) error
	UnRegistry(serName string, serID string) error
}

type Server struct {
	ID       string             `json:"id"`
	Addr     string             `json:"addr"`
	Protocol transport.Protocol `json:"protocol"`
	Weights  float64            `json:"weights"`
}
