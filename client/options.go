package client

import (
	"log"
	"time"

	"github.com/dollarkillerx/light/codes"
	"github.com/dollarkillerx/light/cryptology"
	"github.com/dollarkillerx/light/discovery"
	"github.com/dollarkillerx/light/load_banlancing"
	"github.com/dollarkillerx/light/transport"
)

type Options struct {
	Discovery         discovery.Discovery
	protocol          transport.Protocol
	loadBalancing     load_banlancing.LoadBalancing
	serializationType codes.SerializationType
	compressorType    codes.CompressorType

	pool         int
	cryptology   cryptology.Cryptology
	aesKey       []byte
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func defaultOptions() *Options {
	return &Options{
		pool:              25,
		protocol:          transport.KCP,
		serializationType: codes.MsgPack,
		compressorType:    codes.Snappy,
		loadBalancing:     load_banlancing.NewPolling(),
		cryptology:        cryptology.AES,
		aesKey:            []byte("58a95a8f804b49e686f651a0d3f6e631"),
		writeTimeout:      time.Minute,
		readTimeout:       time.Minute,
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

func SetAESCryptology(key []byte) Option {
	if len(key) != 32 && len(key) != 16 {
		log.Fatalln("AES KEY LEN == 32 OR == 16")
	}

	return func(options *Options) {
		options.cryptology = cryptology.AES
		options.aesKey = key
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
