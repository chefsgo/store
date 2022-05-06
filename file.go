package file

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chefsgo/chef"
)

type (
	filed struct {
		base string
		hash string
		tttt string
		size int64

		file string
		name string
		code string
	}

	File interface {
		Base() string
		Hash() string
		Type() string
		Size() int64

		File() string
		Name() string
		Code() string
	}
	Files []File
)

func (sf *filed) Base() string {
	return sf.base
}
func (sf *filed) Hash() string {
	return sf.hash
}
func (sf *filed) Type() string {
	return sf.tttt
}
func (sf *filed) Size() int64 {
	return sf.size
}
func (sf *filed) File() string {
	return sf.file
}
func (sf *filed) Name() string {
	return sf.name
}
func (sf *filed) Code() string {
	return sf.code
}

// func NewFile(base, hash, filepath string, size int64) File {
// 	file := &filed{}

// 	file.base = base
// 	file.hash = hash
// 	file.path = filepath
// 	file.name = path.Base(file.path)
// 	file.tttt = util.Extension(file.name)
// 	file.size = size

// 	file.code = encode(file)

// 	return file
// }

// func  StatFile(file string) (Map, error) {
// 	stat, err := os.Stat(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hash := util.Sha1File(file)
// 	if hash == "" {
// 		return nil, errors.New("hash error")
// 	}
// 	filename := stat.Name()
// 	extension := util.Extension(file)
// 	mimetype := chef.Mimetype(extension)
// 	length := stat.Size()

// 	return Map{
// 		"hash": hash,
// 		"name": filename,
// 		"type": extension,
// 		"mime": mimetype,
// 		"size": length,
// 		"file": file,
// 	}, nil
// }

//文件编解码
//fileConfig可以设置加解密方式
func encode(info *filed) string {
	code := fmt.Sprintf("%s\t%s\t%s\t%d", info.Base(), info.Hash(), info.Type(), info.Size())
	if val, err := chef.EncryptTEXT(code); err == nil {
		return val
	}
	return ""
}

func decode(code string) *filed {
	val, err := chef.DecryptTEXT(code)
	if err != nil {
		return nil
	}

	args := strings.Split(fmt.Sprintf("%v", val), "\t")
	if len(args) != 4 {
		return nil
	}

	info := &filed{}
	info.base = args[0]
	info.hash = args[1]
	info.tttt = args[2]
	if vv, err := strconv.ParseInt(args[3], 10, 64); err == nil {
		info.size = vv
	}

	return info
}
