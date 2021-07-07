package server

type Server struct {
	serviceMap map[string]*service
	options    *Options
}

func NewServer() *Server {
	return &Server{
		serviceMap: map[string]*service{},
		options:    DefaultOptions(),
	}
}

func (s *Server) Run(options ...Option) error {
	for _, fn := range options {
		fn(s.options)
	}

	return nil
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
