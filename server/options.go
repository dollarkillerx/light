package server

import (
	"context"
	"net"
	"time"
)

type Options struct {
	Protocol Protocol
	UseHttp  bool
	Uri      string
	nl       net.Listener
	ctx      context.Context
	options  map[string]interface{} // 零散配置

	readTimeout     time.Duration
	writeTimeout    time.Duration
	processChanSize int
}

type Protocol string

const (
	TCP  Protocol = "tcp"
	KCP  Protocol = "kcp"
	MQTT Protocol = "mqtt"
)

func (p Protocol) String() string {
	return string(p)
}

func DefaultOptions() *Options {
	return &Options{
		Protocol:     KCP, // default KCP
		Uri:          "0.0.0.0:8397",
		UseHttp:      false,
		readTimeout:  time.Second * 30,
		writeTimeout: time.Second * 30,
		ctx:          context.Background(),
		options: map[string]interface{}{
			"TCPKeepAlivePeriod": time.Minute * 3,
		},
		processChanSize: 1000,
	}
}

type Option func(options *Options)

func UseTCP(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = TCP
	}
}

func UseKCP(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = KCP
	}
}

func UseMQTT(host string) Option {
	return func(options *Options) {
		options.Uri = host
		options.Protocol = MQTT
	}
}

func UseHTTP() Option {
	return func(options *Options) {
		options.Protocol = TCP
		options.UseHttp = true
	}
}

func SetTimeout(readTimeout time.Duration, writeTimeout time.Duration) Option {
	return func(options *Options) {
		options.readTimeout = readTimeout
		options.writeTimeout = writeTimeout
	}
}

func SetContext(ctx context.Context) Option {
	return func(options *Options) {
		options.ctx = ctx
	}
}
