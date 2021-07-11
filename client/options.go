package client

import (
	"log"
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
	aesKey       []byte
	writeTimeout time.Duration
	readTimeout  time.Duration
	heartBeat    time.Duration
	Trace        bool
}

func defaultOptions() *Options {
	return &Options{
		pool:              25,
		serializationType: codes.MsgPack,
		compressorType:    codes.Snappy,
		loadBalancing:     load_banlancing.NewPolling(),
		cryptology:        cryptology.AES,
		aesKey:            []byte("58a95a8f804b49e686f651a0d3f6e631"),
		writeTimeout:      time.Minute,
		readTimeout:       time.Minute,
		heartBeat:         time.Minute * 2,
		Trace:             false,
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
