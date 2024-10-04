package config

import (
	"time"
)

type Configurator interface {
	GetHttpServerName() string
	GetHttpServerPort() string
	GetHttpServerShutDownTimeout() time.Duration
	GetHttpServerRequestTimeout() time.Duration
}
