// +build wasm,electron

package electron

import (
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
	"syscall/js"
)

var electron js.Value

func init() {
	electron = nodejs.Require("electron")
}

func Get() js.Value {
	return electron
}
