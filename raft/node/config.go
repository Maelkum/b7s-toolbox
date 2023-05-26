package node

import (
	"github.com/hashicorp/raft"
)

// Option can be used to set configuration options.
type Option func(cfg *Config)

// Config options for the Node.
type Config struct {
	Bootstrap bool

	// Used to store and retrieve logs.
	LogStore raft.LogStore

	// Stable storage for key configurations.
	StableStore raft.StableStore

	// Snapshot storage and retrieval.
	SnapshotStore raft.SnapshotStore
}

// Should the node bootstrap a Raft cluster.
func BootstrapCluster(should bool) Option {
	return func(cfg *Config) {
		cfg.Bootstrap = should
	}
}

// Specify log store the node should use.
func WithLogStore(logstore raft.LogStore) Option {
	return func(cfg *Config) {
		cfg.LogStore = logstore
	}
}

// Specify stable store the node should use.
func WithStableStore(store raft.StableStore) Option {
	return func(cfg *Config) {
		cfg.StableStore = store
	}
}

func WithSnapshotStore(store raft.SnapshotStore) Option {
	return func(cfg *Config) {
		cfg.SnapshotStore = store
	}
}
