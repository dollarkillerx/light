package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/dollarkillerx/light/transport"
)

type Server struct {
	serviceMap map[string]*service
	options    *Options
}

func NewServer() *Server {
	return &Server{
		serviceMap: map[string]*service{},
		options:    defaultOptions(),
	}
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
		s.options.nl, err = transport.Transport.Gen(transport.KCP, s.options.Uri)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("%s not funod", s.options.Protocol))
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

			go s.process(accept)
		}

	}

	return nil
}
