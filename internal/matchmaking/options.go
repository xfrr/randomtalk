package matchmaking

import (
	"github.com/rs/zerolog"
	matchmakingconfig "github.com/xfrr/randomtalk/internal/matchmaking/config"
)

// InitOption defines the signature for functional options.
type InitOption func(*Service)

// WithConfig overrides the default configuration.
func WithConfig(config matchmakingconfig.Config) InitOption {
	return func(svc *Service) {
		svc.config = config
	}
}

// WithVersion overrides the default version.
func WithVersion(version string) InitOption {
	return func(svc *Service) {
		svc.version = version
	}
}

// WithLogger overrides the default logger.
func WithLogger(logger zerolog.Logger) InitOption {
	return func(svc *Service) {
		svc.logger = &logger
	}
}
