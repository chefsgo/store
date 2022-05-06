package file

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
	"path"

	"54hub.com/argo/util"
	. "github.com/chefsgo/base"
)

type (
	Instance struct {
		Name    string
		Config  Config
		Store   Store
		Setting Map

		connect Connect
	}
)

func (this *Instance) Hash(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()
		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
			// return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}

func (this *Instance) File(hash string, file string, size int64) File {
	info := &filed{}

	info.base = this.Name
	info.hash = hash
	info.file = file
	info.name = path.Base(info.file)
	info.tttt = util.Extension(info.name)
	info.size = size

	info.code = encode(info)

	return file
}
