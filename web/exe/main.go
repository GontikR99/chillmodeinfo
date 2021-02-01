// +build wasm

package main

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/electron/application"
	"github.com/GontikR99/chillmodeinfo/internal/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
)

func main() {
	defer func(){
		if err := recover(); err!=nil {
			nodejs.ConsoleLog(fmt.Sprint(err))
			panic(err)
		}
	}()

	appCtx, exitApp := context.WithCancel(context.Background())

	application.OnReady(func() {
		mainWindow := browserwindow.New(browserwindow.NewConf().
			WithWidth(1024).
			WithHeight(600).
			WithShow(false))
		mainWindow.Once("ready-to-show", func() {
			mainWindow.RemoveMenu()
			mainWindow.Show()
		})
		mainWindow.Once("closed", exitApp)
		mainWindow.LoadFile("index.html")

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
	})

	application.OnWindowAllClosed(exitApp)

	<-appCtx.Done()
	application.Quit()

	<-context.Background().Done()
}
