package store

import (
	"path"
	"sync"
	"time"

	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
)

func init() {
	chef.Register(NAME, module)
}

var (
	module = &Module{
		configs:   make(map[string]Config, 0),
		drivers:   make(map[string]Driver, 0),
		instances: make(map[string]Instance, 0),
	}
)

type (
	Module struct {
		mutex sync.Mutex

		connected, initialized, launched bool

		configs map[string]Config
		drivers map[string]Driver

		instances map[string]Instance

		weights  map[string]int
		hashring *util.HashRing
	}

	Config struct {
		Driver string
		Weight int
		Expiry time.Duration

		Cache string

		Setting Map
	}
	Instance struct {
		name    string
		config  Config
		connect Connect
	}
)

func (this *Module) UploadTo(base string, path string, metadata Map) (File, Files, error) {
	if inst, ok := this.instances[base]; ok {
		return inst.connect.Upload(path, metadata)
	}
	return nil, nil, errInvalidStoreConnection
}

func (this *Module) Upload(path string, metadata Map) (File, Files, error) {
	//这里自动分配一个存储
	base := this.hashring.Locate(path)
	return this.UploadTo(base, path, metadata)
}

func (this *Module) Download(code string) (File, error) {
	file := decode(code)

	if file == nil {
		return nil, errInvalidStoreConnection
	}

	if inst, ok := this.instances[file.Base()]; ok {
		filepath, err := inst.connect.Download(file)
		if err != nil {
			return nil, err
		}

		file.path = filepath
		file.name = path.Base(file.path)

		return file, nil

	}
	return nil, errInvalidStoreConnection
}

func (this *Module) Remove(code string) error {
	file := decode(code)
	if file == nil {
		return errInvalidStoreConnection
	}

	if inst, ok := this.instances[file.Base()]; ok {
		return inst.connect.Remove(file)
	}
	return errInvalidStoreConnection
}
