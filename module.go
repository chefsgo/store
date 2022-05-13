package store

import (
	"sync"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
)

func init() {
	chef.Register(module)
}

var (
	module = &Module{
		configs:   make(Configs, 0),
		drivers:   make(map[string]Driver, 0),
		instances: make(map[string]Instance, 0),
	}
)

type (
	Module struct {
		mutex sync.Mutex

		connected, initialized, launched bool

		configs Configs
		drivers map[string]Driver

		instances map[string]Instance

		weights  map[string]int
		hashring *util.HashRing
	}

	Configs map[string]Config
	Config  struct {
		Driver  string
		Weight  int
		Setting Map
	}
)

func (module *Module) Driver(name string, driver Driver, override bool) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if override {
		module.drivers[name] = driver
	} else {
		if module.drivers[name] == nil {
			module.drivers[name] = driver
		}
	}
}

func (this *Module) Config(name string, config Config, override bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if name == "" {
		name = chef.DEFAULT
	}

	if override {
		this.configs[name] = config
	} else {
		if _, ok := this.configs[name]; ok == false {
			this.configs[name] = config
		}
	}
}
func (this *Module) Configs(config Configs, override bool) {
	for key, val := range config {
		this.Config(key, val, override)
	}
}
