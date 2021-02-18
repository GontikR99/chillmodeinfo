// +build wasm

package jsbinding

import (
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"syscall/js"
)

func MakeArrayBuffer(data []byte) js.Value {
	buffer := js.Global().Get("ArrayBuffer").New(len(data))
	view := js.Global().Get("Uint8Array").New(buffer)
	js.CopyBytesToJS(view, data)
	return buffer
}

func ReadArrayBuffer(buffer js.Value) (data []byte) {
	defer func() {
		if r:=recover(); r!=nil {
			console.Log(r)
			data=nil
			return
		}
	}()

	view := js.Global().Get("Uint8Array").New(buffer)
	length := view.Get("byteLength").Int()
	data = make([]byte, length)
	js.CopyBytesToGo(data, view)
	return data
}