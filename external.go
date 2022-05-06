package store

import (
	"errors"
	"fmt"

	. "github.com/chefsgo/base"
)

func Upload(from Any, metadatas ...Map) (File, Files, error) {
	path := ""
	switch vv := from.(type) {
	case string:
		path = vv
	case Map:
		if file, ok := vv["file"].(string); ok {
			path = file
		} else {
			return nil, nil, errors.New("invalid target")
		}
	default:
		path = fmt.Sprintf("%v", vv)
	}

	var metadata Map
	if len(metadatas) > 0 {
		metadata = metadatas[0]
	}

	return module.Upload(path, metadata)
}

func UploadTo(base string, from Any, metadatas ...Map) (File, Files, error) {
	path := ""
	switch vv := from.(type) {
	case string:
		path = vv
	case Map:
		if file, ok := vv["file"].(string); ok {
			path = file
		} else {
			return nil, nil, errors.New("invalid target")
		}
	default:
		path = fmt.Sprintf("%v", vv)
	}

	var metadata Map
	if len(metadatas) > 0 {
		metadata = metadatas[0]
	}

	return module.UploadTo(base, path, metadata)
}

// func UploadFile(path Any, metadatas ...Map) (File, error) {
// 	file, _, err := Upload(path, metadatas...)
// 	return file, err
// }
// func UploadPath(path Any, metadatas ...Map) (Files, error) {
// 	_, files, err := Upload(path, metadatas...)
// 	return files, err
// }

func Download(code string) (File, error) {
	return module.Download(code)
}
func Remove(code string) error {
	return module.Remove(code)
}

// return mFile.Download(file)
// }
// func Remove(code string) error {
// 	file := mFile.Decode(code)
// 	if file == nil {
// 		return errors.New("无效数据")
// 	}
// 	if file.file != "" {
// 		return mStore.Remove(file)
// 	}
// 	return mFile.Remove(file)
// }
