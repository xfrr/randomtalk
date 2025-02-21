package eventstore

import "errors"

type PersistenceEngine string

const (
	PersistenceEngineMongoDB PersistenceEngine = "mongodb"
	PersistenceEngineNATS    PersistenceEngine = "nats"
)

var (
	// ErrInvalidPersistenceEngine is returned when the provided persistence engine
	// is not supported.
	ErrInvalidPersistenceEngine = errors.New("invalid persistence engine")

	// ErrInvalidConfig is returned when the provided config is invalid
	// for the given persistence engine.
	ErrInvalidConfig = errors.New("invalid config")
)
