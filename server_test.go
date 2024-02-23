package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	server, err := NewServer()
	assert.NoError(t, err)
	assert.NotNil(t, server)
}

func TestRegisterRouters(t *testing.T) {
	server, _ := NewServer()
	rr := NewRouters()
	rr.AddRouter("/test", map[string]HandlerFunc{
		http.MethodGet: func(c Context) error {
			return c.String(http.StatusOK, "test passed")
		},
	})

	_ = server.RegisterRouters(ROOT, rr)

	e := server.GetEcho()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, rr.GetRouters("/test")[0].Methods[http.MethodGet](c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test passed", rec.Body.String())
	}
}

func TestStartAndShutdown(t *testing.T) {
	server, _ := NewServer()

	go server.Start()
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	assert.NoError(t, server.Shutdown(ctx))
}

func TestServerClose(t *testing.T) {
	server, _ := NewServer()

	go server.Start()
	time.Sleep(1 * time.Second)

	assert.NoError(t, server.Close())
}

func TestGetEcho(t *testing.T) {
	server, _ := NewServer()

	assert.IsType(t, &echo.Echo{}, server.GetEcho())
}

func TestGracefulShutdown(t *testing.T) {
	server, _ := NewServer()

	go server.Start()
	time.Sleep(1 * time.Second)

	assert.NoError(t, server.gracefulShutdown())
}
