package client

import (
	"runtime"
	"time"

	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/load_banlancing"
)

type Options struct {
	Discovery         discovery.Discovery
	loadBalancing     load_banlancing.LoadBalancing
	serializationType codes.SerializationType
	compressorType    codes.CompressorType

	pool         int
	cryptology   cryptology.Cryptology
	rsaPublicKey []byte
	writeTimeout time.Duration
	readTimeout  time.Duration
	heartBeat    time.Duration
	Trace        bool
	AUTH         string
}

func defaultOptions() *Options {
	defaultPoolSize := runtime.NumCPU() * 4
	if defaultPoolSize < 20 {
		defaultPoolSize = 20
	}

	return &Options{
		pool:              defaultPoolSize,
		serializationType: codes.MsgPack,
		compressorType:    codes.Snappy,
		loadBalancing:     load_banlancing.NewPolling(),
		cryptology:        cryptology.AES,
		rsaPublicKey: []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDWviNW8C1f+cjy8KF0qT93AA1q
lbQTXPKO4qm34bf6UnSpXgemm1zTEgcPu5Ifka2GgTEgeUMD//iwxr3BTNYA0ARc
soVSN53vklXJqRL3xMWNUFg/2bsAZn5Irlw1xRZfzFzqyCDk5JvUCejvHjvjQwOH
YGHsCfV0pvxPlwFq4wIDAQAB
-----END PUBLIC KEY-----`),
		writeTimeout: time.Minute,
		readTimeout:  time.Minute * 3,
		heartBeat:    time.Minute,
		Trace:        false,
		AUTH:         "",
	}
}

type Option func(options *Options)

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

func SetRASPublicKey(key []byte) Option {
	return func(options *Options) {
		options.cryptology = cryptology.AES
		options.rsaPublicKey = key
	}
}

func SetSerialization(serialization codes.SerializationType) Option {
	return func(options *Options) {
		options.serializationType = serialization
	}
}

func SetCompressor(compressor codes.CompressorType) Option {
	return func(options *Options) {
		options.compressorType = compressor
	}
}

func SetTimeOut(writeTimeout time.Duration, readTimeout time.Duration) Option {
	return func(options *Options) {
		options.writeTimeout = writeTimeout
		options.readTimeout = readTimeout
	}
}

func SetHeartBeat(heartBeat time.Duration) Option {
	return func(options *Options) {
		if heartBeat > 0 {
			options.heartBeat = heartBeat
		}
	}
}

func Trance() Option {
	return func(options *Options) {
		options.Trace = true
	}
}

func SetAUTH(token string) Option {
	return func(options *Options) {
		options.AUTH = token
	}
}
