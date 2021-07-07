package server

import "time"

type Options struct {
	Protocol Protocol
	UseHttp  bool

	readTimeout  time.Duration
	writeTimeout time.Duration
}

type Protocol int

const (
	TCP Protocol = iota
	KCP
	MQTT
)

func DefaultOptions() *Options {
	return &Options{
		Protocol:     KCP, // default KCP
		UseHttp:      false,
		readTimeout:  time.Second * 30,
		writeTimeout: time.Second * 30,
	}
}

type Option func(options *Options)

func UseTCP(host string) Option {
	return func(options *Options) {
		options.Protocol = TCP
	}
}

func UseKCP(host string) Option {
	return func(options *Options) {
		options.Protocol = KCP
	}
}

func UseMQTT(host string) Option {
	return func(options *Options) {
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
