package file

import (
	"sync"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
)

func init() {
	chef.Register(NAME, module)
}

var (
	module = &Module{
		Config:    Config{},
		stores:    make(Stores, 0),
		drivers:   make(map[string]Driver, 0),
		instances: make(map[string]Instance, 0),
	}
)

type (
	Module struct {
		mutex sync.Mutex

		connected, initialized, launched bool

		config Config
		stores Stores

		drivers map[string]Driver

		instances map[string]Instance

		weights  map[string]int
		hashring *util.HashRing
	}

	Config struct {
		Hash string
	}

	Stores map[string]Store
	Store  struct {
		Driver  string
		Weight  int
		Setting Map
	}
)

func (module *Module) Driver(name string, driver Driver, override bool) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if driver == nil {
		panic("Invalid file driver: " + name)
	}

	if override {
		module.drivers[name] = driver
	} else {
		if module.drivers[name] == nil {
			module.drivers[name] = driver
		}
	}
}

func (this *Module) Config(config Config, override bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.config = config
}

func (this *Module) Store(name string, config Store, override bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if override {
		this.stores[name] = config
	} else {
		if _, ok := this.stores[name]; ok == false {
			this.stores[name] = config
		}
	}
}
func (this *Module) Stores(config Stores, override bool) {
	for key, val := range config {
		this.Store(key, val, override)
	}
}
