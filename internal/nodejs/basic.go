// +build wasm

package nodejs

import "syscall/js"

func Require(path string) js.Value {
	return js.Global().Call("require", path)
}

func ConsoleLog(value interface{}) {
	js.Global().Get("console").Call("log", value)
}
