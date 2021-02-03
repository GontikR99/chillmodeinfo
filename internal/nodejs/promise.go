// +build wasm,electron

package nodejs

import (
	"syscall/js"
)

func Promise(prom js.Value) (<-chan js.Value, <-chan js.Value) {
	successChannel:=make(chan js.Value)
	errorChannel:=make(chan js.Value)
	successCallback:=new(js.Func)
	errorCallback:=new(js.Func)

	*successCallback=js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		successChannel <- args[0]
		successCallback.Release()
		errorCallback.Release()
		return nil
	})
	*errorCallback=js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		errorChannel <- args[0]
		successCallback.Release()
		errorCallback.Release()
		return nil
	})
	prom2 := prom.Call("then", *successCallback)
	prom2.Call("catch", *errorCallback)
	return successChannel, errorChannel
}