package file

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
