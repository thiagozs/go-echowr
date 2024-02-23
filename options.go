package server

type Options func(s *ServerParams) error

type ServerParams struct {
	Port string
	Host string
}

func newServerParams(opts ...Options) (*ServerParams, error) {
	s := &ServerParams{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}
func WithPort(port string) Options {
	return func(s *ServerParams) error {
		s.Port = port
		return nil
	}
}

func WithHost(host string) Options {
	return func(s *ServerParams) error {
		s.Host = host
		return nil
	}
}

// getters and setters ------

func (s *ServerParams) GetPort() string {
	return s.Port
}

func (s *ServerParams) GetHost() string {
	return s.Host
}

func (s *ServerParams) SetPort(port string) {
	s.Port = port
}

func (s *ServerParams) SetHost(host string) {
	s.Host = host
}
