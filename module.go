package store

import (
	"errors"
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
		Driver  string
		Weight  int
		Setting Map
	}
	Instance struct {
		name    string
		config  Config
		connect Connect
	}
)

//上传文件，是不是随便选一个库，还是选第一个库
func (this *Module) Upload(base string, path string, metadata Map) (File, Files, error) {
	if inst, ok := this.instances[base]; ok {
		return inst.connect.Upload(path, metadata)
	}

	return nil, nil, errInvalidStoreConnection
}

//下载文件，集成file和store
func (this *Module) Download(file File) (string, error) {
	if file.Store() != "" {
		if conn, ok := this.connects[file.Store()]; ok {
			return conn.Download(file)
		}
		return "", errors.New("无效存储")
	}

	//转给文件
	return mFile.Download(file)
}

func (this *Module) Remove(file File) error {
	if file.Store() != "" {
		if conn, ok := this.connects[file.Store()]; ok {
			return conn.Remove(file)
		}
		return errors.New("无效存储")
	}

	//转给文件
	return mFile.Remove(file)
}

// func (this *Module) Read(key string) (Any, error) {
// 	locate := this.hashring.Locate(key)

// 	if inst, ok := this.instances[locate]; ok {
// 		key := inst.config.Prefix + key //加前缀
// 		return inst.connect.Read(key)
// 	}

// 	return nil, errInvalidStoreConnection
// }
