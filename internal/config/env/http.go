package env

import (
	"errors"
	"net"
	"os"

	"github.com/vadskev/go_final_project/internal/config"
)

const (
	envHostName = "TODO_HOST"
	envPortName = "TODO_PORT"
)

var _ config.HTTPConfig = (*httpConfig)(nil)

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(envHostName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(envPortName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
