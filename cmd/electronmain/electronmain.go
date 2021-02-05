// +build wasm

package main

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/commands"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/application"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
	"time"
)

func main() {
	defer func(){
		if err := recover(); err!=nil {
			console.Log(err)
			panic(err)
		}
	}()

	settings.DefaultSetting(settings.EverQuestDirectory, "C:\\Users\\Public\\Daybreak Game Company\\Installed Games\\EverQuest")
	eqfiles.RestartLogScans()
	commands.WatchLogs()

	appCtx, exitApp := context.WithCancel(context.Background())

	application.OnReady(func() {
		go func() {
			<- time.After(1000*time.Millisecond)
			mainWindow := browserwindow.New(&browserwindow.Conf{
				Width:  1600,
				Height: 800,
				Show:   false,
				WebPreferences: &browserwindow.WebPreferences{
					Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
					NodeIntegration:  false,
					ContextIsolation: true,
				},
			})
			mainWindow.OnClosed(exitApp)
			mainWindow.ServeRPC(exerpcs.NewServer())

			mainWindow.Once("ready-to-show", func() {
				mainWindow.RemoveMenu()
				mainWindow.Show()
			})
			mainWindow.LoadFile(path.Join(application.GetAppPath(), "src/index.html"))
		}()
	})

	application.OnWindowAllClosed(exitApp)

	<-appCtx.Done()
	application.Quit()

	<-context.Background().Done()
}
