// +build wasm

package nodejs

import "syscall/js"

var requireRoot=js.Undefined()

func init() {
	if !js.Global().Get("require").IsUndefined() {
		requireRoot =js.Global()
	}
}

func Require(path string) js.Value {
	if requireRoot.IsUndefined() {
		return js.Undefined()
	} else {
		return requireRoot.Call("require", path)
	}
}

