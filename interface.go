package store

import (
	. "github.com/chefsgo/base"
	"github.com/chefsgo/chef"
	"github.com/chefsgo/util"
	//
)

func (this *Module) Register(name string, value Any, override bool) {
	switch config := value.(type) {
	case Driver:
		this.Driver(name, config, override)
	case Config:
		this.Config(name, config, override)
	case Configs:
		this.Configs(config, override)
	}
}

func (this *Module) configure(name string, config Map) {
	cfg := Config{
		Driver: chef.DEFAULT, Weight: 1,
	}
	//如果已经存在了，用现成的改写
	if vv, ok := this.configs[name]; ok {
		cfg = vv
	}

	if driver, ok := config["driver"].(string); ok {
		cfg.Driver = driver
	}

	//分配权重
	if weight, ok := config["weight"].(int); ok {
		cfg.Weight = weight
	}
	if weight, ok := config["weight"].(int64); ok {
		cfg.Weight = int(weight)
	}
	if weight, ok := config["weight"].(float64); ok {
		cfg.Weight = int(weight)
	}

	if setting, ok := config["setting"].(Map); ok {
		cfg.Setting = setting
	}

	//保存配置
	this.configs[name] = cfg
}
func (this *Module) Configure(global Map) {
	var config Map
	if vvv, ok := global["store"].(Map); ok {
		config = vvv
	}
	if config == nil {
		return
	}

	//记录上一层的配置，如果有的话
	rootConfig := Map{}

	for key, val := range config {
		if conf, ok := val.(Map); ok {
			this.configure(key, conf)
		} else {
			rootConfig[key] = val
		}
	}

	if len(rootConfig) > 0 {
		this.configure(chef.DEFAULT, rootConfig)
	}
}
func (this *Module) Initialize() {
	if this.initialized {
		return
	}

	// 如果没有配置任何连接时，默认一个
	if len(this.configs) == 0 {
		this.configs[chef.DEFAULT] = Config{
			Driver: chef.DEFAULT, Weight: 1,
		}
	} else {
		for key, config := range this.configs {
			if config.Weight == 0 {
				config.Weight = 1
			}
			this.configs[key] = config
		}

	}

	this.initialized = true
}
func (this *Module) Connect() {
	if this.connected {
		return
	}

	//记录要参与分布的连接和权重
	weights := make(map[string]int)

	for name, config := range this.configs {
		driver, ok := this.drivers[config.Driver]
		if ok == false {
			panic("Invalid store driver: " + config.Driver)
		}

		//实例
		inst := Instance{
			name, config, config.Setting, nil,
		}

		// 建立连接
		connect, err := driver.Connect(inst)
		if err != nil {
			panic("Failed to connect to file: " + err.Error())
		}

		// 打开连接
		err = connect.Open()
		if err != nil {
			panic("Failed to open file connect: " + err.Error())
		}

		//记住连接
		inst.connect = connect

		//保存连接
		this.instances[name] = inst

		//只有设置了权重的才参与分布
		if config.Weight > 0 {
			weights[name] = config.Weight
		}
	}

	//hashring分片
	this.weights = weights
	this.hashring = util.NewHashRing(weights)

	this.connected = true
}
func (this *Module) Launch() {
	if this.launched {
		return
	}

	this.launched = true
}
func (this *Module) Terminate() {
	for _, ins := range this.instances {
		ins.connect.Close()
	}

	this.launched = false
	this.connected = false
	this.initialized = false
}
