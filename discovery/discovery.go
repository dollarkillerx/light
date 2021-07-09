package discovery

import "github.com/dollarkillerx/light/transport"

// Discovery 服务发现
type Discovery interface {
	Discovery(serName string) ([]*Server, error)
}

type Server struct {
	Addr     string             `json:"addr"`
	Protocol transport.Protocol `json:"protocol"`
}
