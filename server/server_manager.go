package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dollarkillerx/light"
	"github.com/dollarkillerx/light/transport"
	"github.com/dollarkillerx/light/utils"
)

type Server struct {
	serviceMap map[string]*service
	options    *Options

	beforeMiddleware     []MiddlewareFunc
	afterMiddleware      []MiddlewareFunc
	beforeMiddlewarePath map[string][]MiddlewareFunc
	afterMiddlewarePath  map[string][]MiddlewareFunc
}

func NewServer() *Server {
	return &Server{
		serviceMap: map[string]*service{},
		options:    defaultOptions(),

		beforeMiddleware:     []MiddlewareFunc{},
		afterMiddleware:      []MiddlewareFunc{},
		beforeMiddlewarePath: map[string][]MiddlewareFunc{},
		afterMiddlewarePath:  map[string][]MiddlewareFunc{},
	}
}

type MiddlewareFunc func(ctx *light.Context, request interface{}, response interface{}) error

func (s *Server) Before(beforeMiddleware ...MiddlewareFunc) {
	if len(beforeMiddleware) <= 0 {
		return
	}
	s.beforeMiddleware = append(s.beforeMiddleware, beforeMiddleware...)
}

func (s *Server) After(afterMiddleware ...MiddlewareFunc) {
	if len(afterMiddleware) <= 0 {
		return
	}
	s.afterMiddleware = append(s.afterMiddleware, afterMiddleware...)
}

// BeforePath 前置middleware  path: server_name.server_method
func (s *Server) BeforePath(path string, beforeMiddleware ...MiddlewareFunc) {
	if len(beforeMiddleware) <= 0 {
		return
	}
	fn, ex := s.beforeMiddlewarePath[path]
	if !ex {
		fn = make([]MiddlewareFunc, 0)
	}

	fn = append(fn, beforeMiddleware...)
	s.beforeMiddlewarePath[path] = fn
}

// AfterPath 后置middleware  path: server_name.server_method
func (s *Server) AfterPath(path string, afterMiddleware ...MiddlewareFunc) {
	if len(afterMiddleware) <= 0 {
		return
	}
	fn, ex := s.afterMiddlewarePath[path]
	if !ex {
		fn = make([]MiddlewareFunc, 0)
	}

	fn = append(fn, afterMiddleware...)
	s.afterMiddlewarePath[path] = fn
}

func (s *Server) Register(server interface{}) error {
	return s.register(server, "", false)
}

func (s *Server) RegisterName(server interface{}, serverName string) error {
	return s.register(server, serverName, true)
}

func (s *Server) register(server interface{}, serverName string, useName bool) error {
	ser, err := newService(server, serverName, useName)
	if err != nil {
		return err
	}

	s.serviceMap[ser.name] = ser
	return nil
}

func (s *Server) Run(options ...Option) error {
	for _, fn := range options {
		fn(s.options)
	}

	var err error
	switch s.options.Protocol {
	case transport.KCP:
		s.options.nl, err = transport.Transport.Gen(transport.KCP, s.options.Uri)
		if err != nil {
			return err
		}
	case transport.TCP:
		s.options.nl, err = transport.Transport.Gen(transport.TCP, s.options.Uri)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("%s not funod", s.options.Protocol))
	}

	log.Printf("LightRPC: %s  %s \n", s.options.Protocol, s.options.Uri)

	if s.options.Discovery != nil {
		sIdb, err := ioutil.ReadFile("./light.conf")
		if err != nil {
			id, err := utils.DistributedID()
			if err != nil {
				return err
			}
			sIdb = []byte(id)
		}
		// 进行服务注册
		sId := string(sIdb)
		for k := range s.serviceMap {
			err := s.options.Discovery.Registry(k, s.options.registryAddr, s.options.weights, s.options.Protocol, &sId)
			if err != nil {
				return err
			}
			log.Printf("Discovery Registry: %s addr: %s SUCCESS", k, s.options.registryAddr)
		}

		ioutil.WriteFile("./light.conf", sIdb, 00666)
	}
	return s.run()
}

func (s *Server) run() error {
loop:
	for {
		select {
		case <-s.options.ctx.Done():
			break loop
		default:
			accept, err := s.options.nl.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			if s.options.Trace {
				log.Println("connect: ", accept.RemoteAddr())
			}
			go s.process(accept)
		}

	}

	return nil
}
