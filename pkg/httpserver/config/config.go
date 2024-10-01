package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Configurator interface {
	GetHttpServerName() string
	GetHttpServerPort() string
	GetHttpServerShutDownTimeout() time.Duration
	GetHttpServerRequestTimeout() time.Duration
}

type Config struct {
	// HttpServerName is a name of the shared server.
	HttpServerName string `envconfig:"HTTP_SERVER_NAME" default:"http_server"`
	// HttpServerPort is a port for shared server (endpoints like a /probe for k8s).
	HttpServerPort string `envconfig:"HTTP_SERVER_PORT" default:":8000"`
	// HttpServerShutDownTimeout is a duration value before the server will be closed forcefully.
	HttpServerShutDownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
	// HttpServerRequestTimeout is a timeout value for close request forcefully.
	HttpServerRequestTimeout time.Duration `envconfig:"HTTP_SERVER_REQUEST_TIMEOUT" default:"1m"`
}

func Load() (*Config, error) {
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) GetHttpServerName() string {
	return cfg.HttpServerName
}

func (cfg *Config) GetHttpServerPort() string {
	return cfg.HttpServerPort
}

func (cfg *Config) GetHttpServerShutDownTimeout() time.Duration {
	return cfg.HttpServerShutDownTimeout
}

func (cfg *Config) GetHttpServerRequestTimeout() time.Duration {
	return cfg.HttpServerRequestTimeout
}
