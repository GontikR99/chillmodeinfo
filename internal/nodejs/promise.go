// +build wasm,electron

package nodejs

import (
	"syscall/js"
)

// Interface with a promise and return two channels, the first receiving a success value,
// the second receiving an error value
func FromPromise(prom js.Value) (successChannel <-chan []js.Value, errorChannel <-chan []js.Value) {
	successChannelOut:=make(chan []js.Value)
	successChannel=successChannelOut

	errorChannelOut := make(chan []js.Value)
	errorChannel= errorChannelOut
	successCallback:=new(js.Func)
	errorCallback:=new(js.Func)

	*successCallback=js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		successChannelOut <- args
		successCallback.Release()
		errorCallback.Release()
		return nil
	})
	*errorCallback=js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		errorChannelOut <- args
		successCallback.Release()
		errorCallback.Release()
		return nil
	})
	prom2 := prom.Call("then", *successCallback)
	prom2.Call("catch", *errorCallback)
	return successChannel, errorChannel
}