// +build wasm

package electron

import (
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
	"syscall/js"
)

var electron = nodejs.Require("electron")
var app = electron.Get("app")
var rootDir = js.Global().Get("rootDir")

func Get() js.Value {
	return electron
}

func RootDirectory() string {
	return rootDir.String()
}

