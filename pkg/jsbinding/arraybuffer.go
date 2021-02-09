// +build wasm

package jsbinding

import (
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"syscall/js"
)

func MakeArrayBuffer(data []byte) js.Value {
	view := js.Global().Get("Uint8Array").New(len(data))
	for idx, val := range data {
		view.SetIndex(idx, int(val))
	}
	return view
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
	for i:=0;i<length;i++ {
		data[i]=byte(view.Index(i).Int())
	}
	return data
}