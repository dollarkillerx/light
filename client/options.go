package client

import "github.com/dollarkillerx/light/discovery"

type Options struct {
	Discovery discovery.Discovery
}

func defaultOptions() *Options {
	return &Options{}
}

type Option func(options *Options)

func UseSimpleP2pDiscovery(addr string) Option {
	return func(options *Options) {
		options.Discovery = discovery.NewSimplePeerToPeer(addr)
	}
}
