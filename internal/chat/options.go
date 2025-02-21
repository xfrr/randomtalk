package chatcontext

import (
	"github.com/rs/zerolog"
	chatconfig "github.com/xfrr/randomtalk/internal/chat/config"
)

// InitOption defines the signature for functional options.
type InitOption func(*Service)

// ServiceConfig overrides the default configuration.
func ServiceConfig(config chatconfig.Config) InitOption {
	return func(svc *Service) {
		svc.config = config
	}
}

// ServiceVersion overrides the service version.
func ServiceVersion(version string) InitOption {
	return func(svc *Service) {
		svc.version = version
	}
}

// ServiceLogger overrides the default logger.
func ServiceLogger(logger zerolog.Logger) InitOption {
	return func(svc *Service) {
		svc.logger = &logger
	}
}
