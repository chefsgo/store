package store

import (
	. "github.com/chefsgo/base"
)

type (
	// Driver
	Driver interface {
		Connect(Instance) (Connect, error)
	}

	// Health
	Health struct {
		Workload int64
	}

	// Connect
	Connect interface {
		Open() error
		Health() Health
		Close() error

		Upload(path string, metadata Map) (File, Files, error)
		Download(file File) (string, error)
		Remove(file File) error

		// Browse(file File, name string, expiries ...time.Duration) (string, error)
		// Preview(file File, w, h, t int64, expiries ...time.Duration) (string, error)
	}
)
