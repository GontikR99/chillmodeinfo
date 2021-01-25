// +build wasm

package main

import (
	"context"
	"fmt"
	"syscall/js"
)

func consoleLog(items ...interface{}) {
	js.Global().Get("console").Call("log", items...)
}

func require(path string) js.Value {
	return js.Global().Call("require", path)
}

var electron = require("electron")
var browserWindow = electron.Get("BrowserWindow")
var app = electron.Get("app")

var nodePath = require("path")
var nodeRoot = js.Global().Get("rootDir").String()

func nodePathJoin(parts ...string) string {
	var args []interface{}
	for _, part := range parts {
		args = append(args, part)
	}
	return nodePath.Call("join", args...).String()
}

func singleShot(fn func()) js.Func {
	var fnwrap js.Func
	fnwrap=js.FuncOf(func(_ js.Value, _ []js.Value)interface{} {
		fnwrap.Release()
		fn()
		return nil
	})
	return fnwrap
}

func main() {
	defer func(){
		if err := recover(); err!=nil {
			consoleLog(fmt.Sprint(err))
			panic(err)
		}
	}()

	appCtx, exitApp := context.WithCancel(context.Background())

	js.Global().Get("eventBarriers").Get("ready").Call("onSignal", singleShot(func() {
		mainWindow := browserWindow.New(map[string]interface{}{
			"width":  int(1024),
			"height": int(600),
			"show": false,
		})

		mainWindow.Call("once", "ready-to-show", singleShot(func() {
			mainWindow.Call("removeMenu")
			mainWindow.Call("show")
		}))

		mainWindow.Call("on", "closed", singleShot(func() {
			exitApp()
		}))

		mainWindow.Call("loadFile", nodePathJoin(nodeRoot, "index.html"))

		//overlayWindow := browserWindow.New(map[string]interface{} {
		//	"width": int(400),
		//	"height": int(400),
		//	"show": false,
		//	"transparent": true,
		//	"frame": false,
		//})
		//
		//overlayWindow.Call("once", "ready-to-show", singleShot(func() {
		//	overlayWindow.Call("removeMenu")
		//	overlayWindow.Call("setAlwaysOnTop", true)
		//	overlayWindow.Call("show")
		//}))
		//
		//overlayWindow.Call("loadFile", nodePathJoin(nodeRoot, "overlay.html"))
	}))

	app.Call("on", "window-all-closed", singleShot(exitApp))

	<-appCtx.Done()
	app.Call("quit")

	<-context.Background().Done()
}
