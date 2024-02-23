package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerParams(t *testing.T) {
	params, err := newServerParams(WithPort("8080"), WithHost("localhost"))
	assert.NoError(t, err)
	assert.Equal(t, "8080", params.GetPort())
	assert.Equal(t, "localhost", params.GetHost())
}

func TestWithPort(t *testing.T) {
	params := &ServerParams{}
	option := WithPort("8081")
	err := option(params)
	assert.NoError(t, err)
	assert.Equal(t, "8081", params.GetPort())
}

func TestWithHost(t *testing.T) {
	params := &ServerParams{}
	option := WithHost("127.0.0.1")
	err := option(params)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", params.GetHost())
}

func TestGettersAndSetters(t *testing.T) {
	params := &ServerParams{}

	params.SetPort("8082")
	assert.Equal(t, "8082", params.GetPort())

	params.SetHost("example.com")
	assert.Equal(t, "example.com", params.GetHost())
}
