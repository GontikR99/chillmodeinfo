// +build wasm

package console

import (
	"fmt"
	"syscall/js"
)

func LogRaw(values ...interface{}) {
	js.Global().Get("console").Call("log", values...)
}

func Log(values ...interface{}) {
	LogRaw(fmt.Sprint(values...))
}

func Logf(format string, values ...interface{}) {
	LogRaw(fmt.Sprintf(format, values...))
}
