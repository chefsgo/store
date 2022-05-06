package store

import (
	"path"

	. "github.com/chefsgo/base"
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
	info := decode(code)

	if info == nil {
		return nil, errInvalidStoreConnection
	}

	if inst, ok := this.instances[info.Base()]; ok {
		file, err := inst.connect.Download(info)
		if err != nil {
			return nil, err
		}

		info.file = file
		info.name = path.Base(info.file)

		return info, nil

	}
	return nil, errInvalidStoreConnection
}

func (this *Module) Remove(code string) error {
	info := decode(code)
	if info == nil {
		return errInvalidStoreConnection
	}

	if inst, ok := this.instances[info.Base()]; ok {
		return inst.connect.Remove(info)
	}
	return errInvalidStoreConnection
}
