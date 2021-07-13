package server

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/transport"
)

type AuthFunc func(ctx *light.Context, token string) error

type Options struct {
	Protocol transport.Protocol
	UseHttp  bool
	Uri      string
	nl       net.Listener
	ctx      context.Context
	options  map[string]interface{} // 零散配置
	Trace    bool

	readTimeout     time.Duration
	writeTimeout    time.Duration
	processChanSize int

	AesKey       []byte
	AuthFunc     AuthFunc
	Discovery    discovery.Discovery
	registryAddr string
	weights      float64
}

func defaultOptions() *Options {
	return &Options{
		Protocol:     transport.TCP, // default TCP
		Uri:          "0.0.0.0:8397",
		UseHttp:      false,
		readTimeout:  time.Minute * 3, // 心跳包 默认 3min
		writeTimeout: time.Second * 30,
		ctx:          context.Background(),
		options: map[string]interface{}{
			"TCPKeepAlivePeriod": time.Minute * 3,
		},
		processChanSize: 1000,
		Trace:           false,
		AesKey:          []byte("58a95a8f804b49e686f651a0d3f6e631"),
	}
}

type Option func(options *Options)

func UseTCP(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = transport.TCP
	}
}

func UseUnix(addr string) Option {
	return func(options *Options) {
		options.Uri = addr
		options.Protocol = transport.UNIX
	}
}

func UseKCP(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = transport.KCP
	}
}

func UseMQTT(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = transport.MQTT
	}
}

func UseHTTP() Option {
	return func(options *Options) {
		options.Protocol = transport.TCP
		options.UseHttp = true
	}
}

func SetTimeout(readTimeout time.Duration, writeTimeout time.Duration) Option {
	return func(options *Options) {
		if readTimeout > time.Minute*3 {
			options.readTimeout = readTimeout
		}
		options.writeTimeout = writeTimeout
	}
}

func SetContext(ctx context.Context) Option {
	return func(options *Options) {
		options.ctx = ctx
	}
}

func Trace() Option {
	return func(options *Options) {
		options.Trace = true
	}
}

func SetAESCryptology(key []byte) Option {
	if len(key) != 32 && len(key) != 16 {
		log.Fatalln("AES KEY LEN == 32 OR == 16")
	}

	return func(options *Options) {
		//options.AesKey = cryptology.AES
		options.AesKey = key
	}
}

func SetAUTH(auth AuthFunc) Option {
	return func(options *Options) {
		options.AuthFunc = auth
	}
}

// SetDiscovery discovery, addr 注册本服务地址
func SetDiscovery(discovery discovery.Discovery, addr string, weights float64) Option {
	return func(options *Options) {
		options.Discovery = discovery
		options.registryAddr = addr
		options.weights = weights
	}
}
