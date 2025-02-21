package xnats

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// NewStreamConfig creates a new default StreamConfig with the given context name
// as the stream name and subjects.
func NewStreamConfig(streamName string, subjects ...string) StreamConfig {
	return StreamConfig{
		streamConfig: jetstream.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		},
	}
}

type StreamConfig struct {
	streamConfig jetstream.StreamConfig
}

func (c StreamConfig) WithStreamName(name string) StreamConfig {
	c.streamConfig.Name = name
	return c
}

func (c StreamConfig) WithSubjects(subjects ...string) StreamConfig {
	c.streamConfig.Subjects = subjects
	return c
}

func (c StreamConfig) WithMaxAge(maxAge time.Duration) StreamConfig {
	c.streamConfig.MaxAge = maxAge
	return c
}

func (c StreamConfig) WithMaxBytes(maxBytes int64) StreamConfig {
	c.streamConfig.MaxBytes = maxBytes
	return c
}

func (c StreamConfig) WithMaxMsgSize(maxMsgSize int32) StreamConfig {
	c.streamConfig.MaxMsgSize = maxMsgSize
	return c
}

func (c StreamConfig) WithMaxMsgs(maxMsgs int64) StreamConfig {
	c.streamConfig.MaxMsgs = maxMsgs
	return c
}

func (c StreamConfig) WithMaxMsgsPerSubject(maxMsgsPerSubject int64) StreamConfig {
	c.streamConfig.MaxMsgsPerSubject = maxMsgsPerSubject
	return c
}

func (c StreamConfig) WithStorageType(storageType jetstream.StorageType) StreamConfig {
	c.streamConfig.Storage = storageType
	return c
}

func (c StreamConfig) WithReplicas(replicas int) StreamConfig {
	c.streamConfig.Replicas = replicas
	return c
}

func (c StreamConfig) WithRetention(retention jetstream.RetentionPolicy) StreamConfig {
	c.streamConfig.Retention = retention
	return c
}

func (c StreamConfig) WithDiscardPolicy(discardPolicy jetstream.DiscardPolicy) StreamConfig {
	c.streamConfig.Discard = discardPolicy
	return c
}

func (c StreamConfig) WithMaxConsumers(maxConsumers int) StreamConfig {
	c.streamConfig.MaxConsumers = maxConsumers
	return c
}

func (c StreamConfig) WithDenyDelete() StreamConfig {
	c.streamConfig.DenyDelete = true
	return c
}

func (c StreamConfig) WithDenyPurge() StreamConfig {
	c.streamConfig.DenyPurge = true
	return c
}
