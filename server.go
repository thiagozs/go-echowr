package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Kind int

const (
	ROOT Kind = iota
	V1
	V2
	V3
	DEV
	API
	DOCS
)

func (k Kind) String() string {
	return [...]string{
		"root",
		"v1",
		"v2",
		"v3",
		"dev",
		"api",
		"docs",
	}[k]
}

type RegisterRouter struct {
	Path    string
	Methods map[string]HandlerFunc
}
type RegisterRouters struct {
	PathFixed string
	Routers   []RegisterRouter
}

func NewRouters() *RegisterRouters {
	return &RegisterRouters{}
}

func (r *RegisterRouters) AddRouter(path string, methods map[string]HandlerFunc) {
	r.Routers = append(r.Routers, RegisterRouter{
		Path:    path,
		Methods: methods,
	})
}

func (r *RegisterRouters) AddRouterFx(params string, methods map[string]HandlerFunc) {
	path := strings.TrimSpace(params)
	if len(path) > 0 {
		path = r.PathFixed + path
	} else {
		path = r.PathFixed
	}

	r.Routers = append(r.Routers, RegisterRouter{
		Path:    path,
		Methods: methods,
	})
}

func (r *RegisterRouters) GeAlltRouters() []RegisterRouter {
	return r.Routers
}

func (r *RegisterRouters) GetRouters(path string) []RegisterRouter {
	var routers []RegisterRouter
	for _, router := range r.Routers {
		if router.Path == path {
			routers = append(routers, router)
		}
	}
	return routers
}

func (r *RegisterRouters) GetRoutersFx() []RegisterRouter {
	var routers []RegisterRouter

	for _, router := range r.Routers {
		if strings.Contains(router.Path, r.PathFixed) {
			routers = append(routers, router)
		}
	}

	return routers
}

func (r *RegisterRouters) SetPathFixed(path string) {
	r.PathFixed = path
}

type Methods map[string]HandlerFunc

type HandlerFunc = echo.HandlerFunc

type MiddlewareFunc = echo.MiddlewareFunc

type Context = echo.Context

type Route = echo.Route

type Server struct {
	port   string
	host   string
	echo   *echo.Echo
	params *ServerParams
}

func NewServer(opts ...Options) (*Server, error) {
	params, err := newServerParams(opts...)
	if err != nil {
		return nil, err
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.HideBanner = true

	return &Server{
		echo:   e,
		port:   params.GetPort(),
		host:   params.GetHost(),
		params: params,
	}, nil
}

func (s *Server) NewContext(req *http.Request, w http.ResponseWriter) Context {
	return s.echo.NewContext(req, w)
}

func (s *Server) RegisterRouters(group Kind, routers *RegisterRouters, middlewares ...MiddlewareFunc) error {
	var grp *echo.Group

	switch group {
	case ROOT:
	case V1, V2, V3, DEV, API, DOCS:
		grp = s.echo.Group(group.String())
	}

	if grp != nil {
		return s.registerRouters(grp, routers, middlewares...)
	}

	return s.registerRouters(s.echo, routers, middlewares...)
}

func (s *Server) registerRouters(engine any, routers *RegisterRouters, middlewares ...MiddlewareFunc) error {
	switch e := engine.(type) {
	case *echo.Group:
		for _, middleware := range middlewares {
			e.Use(middleware)
		}
		for _, methods := range routers.GeAlltRouters() {
			for method, handler := range methods.Methods {
				switch method {
				case http.MethodGet:
					e.GET(methods.Path, handler)
				case http.MethodPost:
					e.POST(methods.Path, handler)
				case http.MethodPut:
					e.PUT(methods.Path, handler)
				case http.MethodDelete:
					e.DELETE(methods.Path, handler)
				case http.MethodPatch:
					e.PATCH(methods.Path, handler)
				case http.MethodHead:
					e.HEAD(methods.Path, handler)
				case http.MethodConnect:
					e.CONNECT(methods.Path, handler)
				case http.MethodOptions:
					e.OPTIONS(methods.Path, handler)
				case http.MethodTrace:
					e.TRACE(methods.Path, handler)
				}
			}
		}

	case *echo.Echo:
		for _, middleware := range middlewares {
			e.Use(middleware)
		}

		for _, methods := range routers.GeAlltRouters() {
			for method, handler := range methods.Methods {
				switch method {
				case http.MethodGet:
					e.GET(methods.Path, handler)
				case http.MethodPost:
					e.POST(methods.Path, handler)
				case http.MethodPut:
					e.PUT(methods.Path, handler)
				case http.MethodDelete:
					e.DELETE(methods.Path, handler)
				case http.MethodPatch:
					e.PATCH(methods.Path, handler)
				case http.MethodHead:
					e.HEAD(methods.Path, handler)
				case http.MethodConnect:
					e.CONNECT(methods.Path, handler)
				case http.MethodOptions:
					e.OPTIONS(methods.Path, handler)
				case http.MethodTrace:
					e.TRACE(methods.Path, handler)
				}
			}
		}

	default:
		return fmt.Errorf("engine type not supported")
	}

	return nil
}

func (s *Server) Start() {
	host := fmt.Sprintf("%s:%s", s.host, s.port)

	if len(s.port) == 0 {
		host = s.host
	}

	go func() {
		if err := s.echo.Start(host); err != nil &&
			err != http.ErrServerClosed {
			s.echo.Logger.Fatal(err)
		}
	}()

}

func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}

func (s *Server) GetRouters() []*Route {
	return s.echo.Routes()
}

func (s *Server) Close() error {
	return s.echo.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *Server) GracefulShutdown() error {
	return s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}
