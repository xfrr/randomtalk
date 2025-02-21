package eventstore

import "time"

// Config is the default configuration for the event store.
// It uses tags to specify the environment variables to read from and their default values.
type Config struct {
	// Environment is the environment the service is running in.
	Environment string `env:"ENVIRONMENT" envDefault:"development"`

	// EventStorePersistenceEngine is the engine to use for the event store.
	// Supported values are 'nats' and 'mongodb'.
	// Defaults to nats.
	EventStorePersistenceEngine string `env:"EVENTSTORE_PERSISTENCE_ENGINE,notEmpty" envDefault:"nats"`
	// EventStoreStreamName is the name of the stream to append events to.
	EventStoreStreamName string `env:"EVENTSTORE_STREAM_NAME" envDefault:"assets-events"`
	// EventStoreStreamMaxAge is the maximum time after which events will be deleted.
	EventStoreStreamMaxAge time.Duration `env:"EVENTSTORE_STREAM_MAX_AGE" envDefault:"8760h"` // 1 year
	// EventStoreStreamReplicas is the number of stream replicas to keep.
	EventStoreStreamReplicas int `env:"EVENTSTORE_STREAM_REPLICAS" envDefault:"1"`

	// EventStoreConnectionDrainTimeout is the maximum amount of time to wait for a connection to drain.
	EventStoreConnectTimeout time.Duration `env:"EVENTSTORE_CONNECT_TIMEOUT" envDefault:"5s"`

	// EventStoreConnectionDrainTimeout is the maximum amount of time to wait for a connection to drain.
	// Drain is the process of waiting for all pending messages to be sent before closing the connection.
	EventStoreConnectionDrainTimeout time.Duration `env:"EVENTSTORE_CONNECT_DRAIN_TIMEOUT" envDefault:"5s"`
	// EventStoreConnectRetryAttempts is the maximum number of times to retry connecting to the event store.
	EventStoreConnectRetryAttempts int `env:"EVENTSTORE_CONNECT_RETRY_LIMIT" envDefault:"3"`
	// EventStoreConnectRetryDelay is the amount of time to wait before retrying to connect to the event store.
	EventStoreConnectRetryDelay time.Duration `env:"EVENTSTORE_CONNECT_RETRY_DELAY" envDefault:"1s"`
	// EventStoreConnectionRetryFactor is the factor to multiply the delay by after each retry.
	EventStoreConnectionRetryFactor time.Duration `env:"EVENTSTORE_CONNECT_RETRY_FACTOR" envDefault:"2s"`

	// EventStoreStreamAppendRetryLimit is the maximum number of times to retry appending to a stream.
	EventStoreStreamAppendRetryLimit int `env:"EVENTSTORE_STREAM_APPEND_RETRY_LIMIT" envDefault:"3"`
	// EventStoreStreamAppendRetryDelay is the amount of time to wait before retrying to append to a stream.
	EventStoreStreamAppendRetryDelay time.Duration `env:"EVENTSTORE_STREAM_APPEND_RETRY_DELAY" envDefault:"1s"`
	// EventStoreStreamAppendRetryFactor is the factor to multiply the delay by after each retry.
	EventStoreStreamAppendRetryFactor time.Duration `env:"EVENTSTORE_STREAM_APPEND_RETRY_FACTOR" envDefault:"2s"`

	// EventStoreMongoDBDatabase is the name of the MongoDB database to use.
	EventStoreMongoDBDatabase string `env:"EVENTSTORE_MONGODB_DATABASE" envDefault:"randomtalk"`
	// EventStoreMongoDBConnectionURI is the URI of the MongoDB server to connect to.
	EventStoreMongoDBConnectionURI string `env:"EVENTSTORE_MONGODB_CONNECT_URI" envDefault:"localhost:27017"`
	// EventStoreMongoDBConnectionUser is the username to use when connecting to the MongoDB server.
	EventStoreMongoDBConnectionUser string `env:"EVENTSTORE_MONGODB_CONNECT_USER,unset"`
	// EventStoreMongoDBConnectionPass is the password to use when connecting to the MongoDB server.
	EventStoreMongoDBConnectionPass string `env:"EVENTSTORE_MONGODB_CONNECT_PASS,unset"`
	// EventStoreMongoDBReadConcern is the read concern to use when reading from the MongoDB server.
	EventStoreMongoDBReadConcern string `env:"EVENTSTORE_MONGODB_READ_CONCERN" envDefault:"majority"`
	// EventStoreMongoDBReadPreference is the read preference to use when reading from the MongoDB server.
	EventStoreMongoDBReadPreference string `env:"EVENTSTORE_MONGODB_READ_PREFERENCE" envDefault:"primary"`
	// EventStoreMongoDBWriteConcern is the write concern to use when writing to the MongoDB server.
	EventStoreMongoDBWriteConcern string `env:"EVENTSTORE_MONGODB_WRITE_CONCERN" envDefault:"majority"`
	// EventStoreMongoDBAuthSource is the authentication source to use when connecting to the MongoDB server.
	EventStoreMongoDBAuthSource string `env:"EVENTSTORE_MONGODB_AUTH_SOURCE" envDefault:"admin"`

	// EventStoreNATSConnectionURI is the URI of the NATS server to connect to.
	EventStoreNATSConnectionURI string `env:"EVENTSTORE_NATS_CONNECT_URI" envDefault:"localhost:4222"`
	// EventStoreNATSConnectionUser is the username to use when connecting to the NATS server.
	EventStoreNATSConnectionUser string `env:"EVENTSTORE_NATS_CONNECT_USER,unset"`
	// EventStoreNATSConnectionPass is the password to use when connecting to the NATS server.
	EventStoreNATSConnectionPass string `env:"EVENTSTORE_NATS_CONNECT_PASS,unset"`
}
