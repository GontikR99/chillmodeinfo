// +build wasm

package application

import (
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"syscall/js"
)

var app=electron.Get().Get("app")

func Quit() {
	app.Call("quit")
}

func GetAppPath() string {
	return app.Call("getAppPath").String()
}

func On(eventName string, handler func(event js.Value)) {
	app.Call("on", eventName, js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		handler(args[0])
		return nil
	}))
}

func OnWindowAllClosed(handler func()) {
	On("window-all-closed", func(event js.Value){
		handler()
	})
}

func OnReady(readyFunc func()) {
	var readyWrapped js.Func
	readyWrapped = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		readyFunc()
		readyWrapped.Release()
		return nil
	})
	js.Global().Get("eventBarriers").Get("ready").Call("onSignal", readyWrapped)
}