// +build wasm

package console

import (
	"fmt"
	"syscall/js"
)

func Log(values ...interface{}) {
	js.Global().Get("console").Call("log", values...)
}

func Logf(format string, values ...interface{}) {
	Log(fmt.Sprintf(format, values...))
}
