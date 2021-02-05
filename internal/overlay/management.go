// +build wasm,electron

package overlay

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/application"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs/path"
	"time"
)

var activeOverlays=make(map[string]browserwindow.BrowserWindow)

// Force a given overlay to close
func Close(fileRoot string) {
	if bw, ok := activeOverlays[fileRoot]; ok {
		end := make(chan struct{})
		bw.OnClosed(func() {
			close(end)
		})
		bw.Destroy()
		<-end
	}
}

// See if a given overlay currently exists
func Lookup(fileRoot string) browserwindow.BrowserWindow {
	bw, _ := activeOverlays[fileRoot]
	return bw
}

const sizingPrefix="windowBounds:"

func GetPreferredSizing(fileRoot string) *electron.Rectangle {
	text, present, err := settings.LookupSetting(sizingPrefix+fileRoot)
	if err!=nil || !present {
		return nil
	}
	result:=new(electron.Rectangle)
	err = json.Unmarshal([]byte(text), result)
	if err!=nil {
		return nil
	}
	return result
}

func SetPreferredSizing(fileRoot string, sizing *electron.Rectangle) {
	data, err := json.Marshal(sizing)
	if err!=nil {
		return
	}
	settings.SetSetting(sizingPrefix+fileRoot, string(data))
}

func ClearPreferredSizing(fileRoot string) {
	settings.ClearSetting(sizingPrefix+fileRoot)
}

func UpdateSizing(title string, fileRoot string) {
	Close(fileRoot)
	psp:=new(*electron.Rectangle)
	*psp = GetPreferredSizing(fileRoot)

	overlayWindow := browserwindow.New(&browserwindow.Conf{
		Title: title,
		Width: 400,
		Height: 400,
		Show: false,
		Transparent: false,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src/preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	activeOverlays[fileRoot]=overlayWindow
	end:=make(chan struct{})
	overlayWindow.Once("ready-to-show", func() {
		if *psp!=nil {
			overlayWindow.SetBounds(*psp)
		}
		overlayWindow.RemoveMenu()
		overlayWindow.Show()
		overlayWindow.SetAlwaysOnTop(true)
		go func() {
			defer func() {recover()}()
			for {
				select {
				case <-end:
					return
				case <-time.After(200*time.Millisecond):
					*psp = overlayWindow.GetBounds()
				}
			}
		}()
	})
	overlayWindow.Once("closed", func() {
		delete(activeOverlays, fileRoot)
		close(end)
	})
	overlayWindow.LoadFile(path.Join(application.GetAppPath(), "src/overlay_position.html"))
	<-end
	if *psp != nil {
		SetPreferredSizing(fileRoot, *psp)
	}
}

func Launch(fileRoot string, interactive bool) browserwindow.BrowserWindow {
	Close(fileRoot)
	prefLoc := GetPreferredSizing(fileRoot)

	overlayWindow := browserwindow.New(&browserwindow.Conf{
		Width: 400,
		Height: 400,
		Show: false,
		Transparent: true,
		Frame: false,
		Resizable: false,
		WebPreferences: &browserwindow.WebPreferences{
			Preload:          path.Join(application.GetAppPath(), "src", "preload.js"),
			NodeIntegration:  false,
			ContextIsolation: true,
		},
	})
	activeOverlays[fileRoot]=overlayWindow
	overlayWindow.Once("ready-to-show", func() {
		if prefLoc!=nil {
			overlayWindow.SetContentBounds(prefLoc)
		}
		overlayWindow.RemoveMenu()
		overlayWindow.ShowInactive()
		overlayWindow.SetAlwaysOnTop(true)
		if !interactive {
			overlayWindow.SetIgnoreMouseEvents(true)
		}
	})
	overlayWindow.OnClosed(func() {
		delete(activeOverlays, fileRoot)
	})
	overlayWindow.LoadFile(path.Join(application.GetAppPath(), "src", fileRoot))

	return overlayWindow
}