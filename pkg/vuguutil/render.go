// +build wasm,web

package vuguutil

type CallbackHandle int

var callbackHandleGen=CallbackHandle(0)
var renderCallbacks=make(map[CallbackHandle]func())

func OnRender(callback func()) CallbackHandle {
	callbackHandleGen++
	renderCallbacks[callbackHandleGen]=callback
	return callbackHandleGen
}

func (ch CallbackHandle) Release() {
	delete(renderCallbacks, ch)
}

func InvokeRenderCallbacks() {
	for _, v := range renderCallbacks {
		v()
	}
}