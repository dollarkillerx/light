package client

import (
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/load_banlancing"
	"github.com/dollarkillerx/light/transport"
)

type Options struct {
	Discovery     discovery.Discovery
	protocol      transport.Protocol
	loadBalancing load_banlancing.LoadBalancing
	pool          int
}

func defaultOptions() *Options {
	return &Options{
		pool:          25,
		protocol:      transport.KCP,
		loadBalancing: load_banlancing.NewPolling(),
	}
}

type Option func(options *Options)

func UseSimpleP2pDiscovery(addr string) Option {
	return func(options *Options) {
		options.Discovery = discovery.NewSimplePeerToPeer(addr)
	}
}

func SetPoolSize(p int) Option {
	if p <= 0 {
		p = 25
	}

	return func(options *Options) {
		options.pool = p
	}
}

func SetLoadBalancing(lb load_banlancing.LoadBalancing) Option {
	return func(options *Options) {
		options.loadBalancing = lb
	}
}
