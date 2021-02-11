// +build wasm,electron

package electron

import (
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"syscall/js"
)

var electron js.Value

func init() {
	electron = nodejs.Require("electron")
}

func JSValue() js.Value {
	return electron
}