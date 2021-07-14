package server

import (
	"context"
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

	MaximumLoad   int64
	RSAPublicKey  []byte
	RSAPrivateKey []byte
	AuthFunc      AuthFunc
	Discovery     discovery.Discovery
	registryAddr  string
	weights       float64
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
		RSAPublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDWviNW8C1f+cjy8KF0qT93AA1q
lbQTXPKO4qm34bf6UnSpXgemm1zTEgcPu5Ifka2GgTEgeUMD//iwxr3BTNYA0ARc
soVSN53vklXJqRL3xMWNUFg/2bsAZn5Irlw1xRZfzFzqyCDk5JvUCejvHjvjQwOH
YGHsCfV0pvxPlwFq4wIDAQAB
-----END PUBLIC KEY-----`),
		RSAPrivateKey: []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDWviNW8C1f+cjy8KF0qT93AA1qlbQTXPKO4qm34bf6UnSpXgem
m1zTEgcPu5Ifka2GgTEgeUMD//iwxr3BTNYA0ARcsoVSN53vklXJqRL3xMWNUFg/
2bsAZn5Irlw1xRZfzFzqyCDk5JvUCejvHjvjQwOHYGHsCfV0pvxPlwFq4wIDAQAB
AoGBAJWh07omjVeNE8rEhZxmuoRPEwor2liLsbCCnEQ3Eh1pC0Vg8e/T3jBtJWJ/
DujUd5d7uiGonVvSJxX2xg5FXe/Xo2bkDL98+mL2MrSBajLR0wB8WcWPpfhvblwJ
n9SEfWjCZdPQQdm6+tGkJaSLFoCDRVkkO4pijaX+S0VWT+kBAkEA3HwftjDw+xHQ
WDkO/pVRU7OVLdBmukaU5itx2SwKoMta7J1MQe0z1Y2gnZhwkL1NeMoYXuul0wMT
7JR2vMikIwJBAPlVO4rbSIbzKL4Zkx8C9SjvJ57kFqU/kRvwg6nVLCNIzaVDphG/
E8v+jo8KoSX7Gyf0xR1xZMcQSbjF2Wd2KkECQEikPG5yQXL2s4Xdhqsp1tmU2Rl3
B+FnT7dlqOS8NeQ0G4jJak5uMB2zw68ogi2tsNCTBOSBDukuomnXoCcik7ECQQDO
HUCQpIALVz4qEHhHnalPQozNVB7IUolBwI0HO3s2W/vsj8TcTMovy+rLouzeufuU
B0tf8Jpv2S4oeh4j4lJBAkEAgykTI0fVKxt9yl4p/xcrLBgs0IBaeC+VolyMBc3H
TT0p2Vmye/dE/9DmugwURhalDEywX7EO0THHr2hg0zhVEQ==
-----END RSA PRIVATE KEY-----`),
		Discovery: &discovery.SimplePeerToPeer{},
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

func SetRSAKey(publicKey []byte, privateKey []byte) Option {
	return func(options *Options) {
		options.RSAPublicKey = publicKey
		options.RSAPrivateKey = privateKey
	}
}

func SetAUTH(auth AuthFunc) Option {
	return func(options *Options) {
		options.AuthFunc = auth
	}
}

// SetDiscovery discovery, addr 注册本服务地址
func SetDiscovery(discovery discovery.Discovery, addr string, weights float64, maximumLoad int64) Option {
	return func(options *Options) {
		options.Discovery = discovery
		options.registryAddr = addr
		options.weights = weights
		options.MaximumLoad = maximumLoad
	}
}
