// +build wasm

package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/commands"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/application"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
)

func main() {
	defer func(){
		if err := recover(); err!=nil {
			console.Log(err)
			panic(err)
		}
	}()

	if _, present, err := settings.LookupSetting(settings.ClientId); err==nil && !present {
		clientId := make([]byte, 32)
		rand.Read(clientId)
		settings.SetSetting(settings.ClientId, hex.EncodeToString(clientId))
	}

	settings.DefaultSetting(settings.EverQuestDirectory, "C:\\Users\\Public\\Daybreak Game Company\\Installed Games\\EverQuest")
	eqspec.RestartLogScans()
	commands.WatchLogs()
	localprofile.StartElectronPoll()

	appCtx, exitApp := context.WithCancel(context.Background())

	application.OnReady(func() {
		go func() {
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
				//mainWindow.RemoveMenu()
				mainWindow.Show()
//				go eqspec.BuildTrie()
			})
			mainWindow.LoadFile(path.Join(application.GetAppPath(), "src/index.html"))
		}()
	})

	application.OnWindowAllClosed(exitApp)

	<-appCtx.Done()
	application.Quit()

	<-context.Background().Done()
}
