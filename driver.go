package store

import (
	. "github.com/chefsgo/base"
)

type (
	// Driver
	Driver interface {
		Connect(name string, config Config) (Connect, error)
	}

	// Health
	Health struct {
		Workload int64
	}

	// Connect
	Connect interface {
		Open() error
		Health() (Health, error)
		Close() error

		Upload(path string, metadata Map) (File, Files, error)
		Download(file File) (string, error)
		Remove(file File) error

		// Browse(file File, name string, expiries ...time.Duration) (string, error)
		// Preview(file File, w, h, t int64, expiries ...time.Duration) (string, error)
	}
)

// Driver 注册驱动
func (module *Module) Driver(name string, driver Driver, override bool) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if driver == nil {
		panic("Invalid store driver: " + name)
	}

	if override {
		module.drivers[name] = driver
	} else {
		if module.drivers[name] == nil {
			module.drivers[name] = driver
		}
	}
}
