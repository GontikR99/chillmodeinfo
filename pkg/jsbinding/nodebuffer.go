// +build wasm,electron

package jsbinding

import "syscall/js"

func BufferOf(data []byte) js.Value {
	view := js.Global().Get("Buffer").Call("alloc", len(data))
	for idx, val := range data {
		view.SetIndex(idx, int(val))
	}
	return view
}