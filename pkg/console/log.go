// +build wasm

package console

import (
	"fmt"
	"syscall/js"
)

func LogRaw(value js.Value) {
	js.Global().Get("console").Call("log", value)
}

func Log(values ...interface{}) {
	LogRaw(js.ValueOf(fmt.Sprint(values...)))
}

func Logf(format string, values ...interface{}) {
	LogRaw(js.ValueOf(fmt.Sprintf(format, values...)))
}
