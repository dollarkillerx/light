package load_banlancing

import "github.com/dollarkillerx/light/discovery"

type LoadBalancing interface {
	InitBalancing(ser []*discovery.Server)
	GetService() (*discovery.Server, error)
}
