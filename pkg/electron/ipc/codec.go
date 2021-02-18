// +build wasm

package ipc

import (
	"github.com/GontikR99/chillmodeinfo/pkg/jsbinding"
	"syscall/js"
)

func Encode(data []byte) js.Value {
	return jsbinding.MakeArrayBuffer(data)
}

func Decode(encoded js.Value) ([]byte,error) {
	return jsbinding.ReadArrayBuffer(encoded), nil
}