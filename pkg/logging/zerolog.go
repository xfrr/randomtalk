package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xfrr/randomtalk/pkg/env"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func NewLogger(serviceName string, environment env.Environment, lvlStr string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	lvl := zerolog.InfoLevel
	if lvlStr != "" {
		lvl, _ = zerolog.ParseLevel(lvlStr)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Str(string(semconv.ServiceNameKey), serviceName).
		Str(string(semconv.DeploymentEnvironmentKey), environment.String()).
		Timestamp().
		Logger().
		Level(lvl)

	log.Logger = logger
	return logger
}

func SetGlobalLogger(logger zerolog.Logger) {
	zerolog.SetGlobalLevel(logger.GetLevel())
	log.Logger = logger
}
