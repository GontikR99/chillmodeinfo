// +build wasm,electron

package fs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"os"
	"syscall/js"
	"time"
)

func ReadDir(path string) ([]os.FileInfo, error) {
	successChan, errorChan := nodejs.FromPromise(fsPromises.Call("readdir", path, map[string]interface{}{"withFileTypes": true}))
	select {
	case errJs:=<-errorChan:
		return nil, errors.New(errJs[0].String())
	case succObj:=<-successChan:
		var files []os.FileInfo
		for i:=0;i<succObj[0].Length();i++ {
			jsEnt := succObj[0].Index(i)
			files = append(files, &jsDirEnt{jsEnt})
		}
		return files, nil
	}
}

type jsDirEnt struct {
	v js.Value
}

func (j *jsDirEnt) Name() string {
	return j.v.Get("name").String()
}

func (j *jsDirEnt) Size() int64 {
	return -1
}

func (j *jsDirEnt) Mode() os.FileMode {
	mode := os.FileMode(0)
	if j.v.Call("isDirectory").Bool() {
		mode |= os.ModeDir
	}
	if j.v.Call("isCharacterDevice").Bool() {
		mode |= os.ModeCharDevice
	}
	if j.v.Call("isFIFO").Bool() {
		mode |= os.ModeNamedPipe
	}
	if j.v.Call("isSocket").Bool() {
		mode |= os.ModeSocket
	}
	if j.v.Call("isSymbolicLink").Bool() {
		mode |= os.ModeSymlink
	}
	if j.v.Call("isBlockDevice").Bool() {
		mode |= os.ModeDevice
	}
	return mode
}

func (j *jsDirEnt) ModTime() time.Time {
	return time.Unix(0,0)
}

func (j *jsDirEnt) IsDir() bool {
	return j.v.Call("isDirectory").Bool()
}

func (j *jsDirEnt) Sys() interface{} {
	panic("implement me")
}
