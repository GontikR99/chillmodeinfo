// +build wasm

package main

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/electron/application"
	"github.com/GontikR99/chillmodeinfo/internal/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcmain"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs/path"
)

func main() {
	defer func(){
		if err := recover(); err!=nil {
			nodejs.ConsoleLog(fmt.Sprint(err))
			panic(err)
		}
	}()

	appCtx, exitApp := context.WithCancel(context.Background())

	pingChan := ipcmain.Listen("pings")
	go func() {
		for {
			msg := <- pingChan
			nodejs.ConsoleLog(string(msg.Content()))
		}
	}()

	application.OnReady(func() {
		mainWindow := browserwindow.New(browserwindow.Conf{
			Width:  1024,
			Height: 600,
			Show:   false,
			WebPreferences: &browserwindow.WebPreferences{
				Preload: path.Join(application.GetAppPath(), "src/preload.js"),
				NodeIntegration: false,
				ContextIsolation: true,
			},
		})
		mainWindow.Once("ready-to-show", func() {
			mainWindow.RemoveMenu()
			mainWindow.Show()
		})
		mainWindow.Once("closed", exitApp)
		mainWindow.LoadFile(path.Join(application.GetAppPath(), "src/index.html"))

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
