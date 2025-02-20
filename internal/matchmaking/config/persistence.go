package matchmakingconfig

import "github.com/kelseyhightower/envconfig"

// UserStoreEngineType is the type of the user store engine.
type UserStoreEngineType string

func (t UserStoreEngineType) String() string {
	return string(t)
}

func (t UserStoreEngineType) IsValid() bool {
	switch t {
	case UserStoreEngineMemory:
		return true
	default:
		return false
	}
}

const (
	// UserStoreEngineMemory is the memory engine.
	UserStoreEngineMemory UserStoreEngineType = "memory"
)

// MatchRepositoryEngineType is the type of the match repository engine.
type MatchRepositoryEngineType string

func (t MatchRepositoryEngineType) String() string {
	return string(t)
}

func (t MatchRepositoryEngineType) IsValid() bool {
	switch t {
	case MatchRepositoryEngineMemory:
		return true
	default:
		return false
	}
}

const (
	// MatchRepositoryEngineMemory is the memory engine.
	MatchRepositoryEngineMemory MatchRepositoryEngineType = "memory"

	// MatchRepositoryEngineNATS is the NATS JetStream engine.
	MatchRepositoryEngineNATS MatchRepositoryEngineType = "nats"
)

// Persistence holds the configuration for the Persistence layer.
type Persistence struct {
	// UserStoreEngine is the engine used for the user store.
	UserStoreEngine UserStoreEngineType `envconfig:"USER_STORE_ENGINE" default:"memory"`

	// MatchRepositoryEngine is the engine used for the match repository.
	MatchRepositoryEngine MatchRepositoryEngineType `envconfig:"MATCH_REPOSITORY_ENGINE" default:"nats"`
}

func mustLoadPersistenceConfig() Persistence {
	var cfg Persistence
	envconfig.MustProcess(envPrefix, &cfg)
	return cfg
}
