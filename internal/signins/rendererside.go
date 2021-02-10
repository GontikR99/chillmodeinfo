// +build wasm,web

package signins

import (
	"bytes"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
)

func init() {
	if electron.IsPresent() {
		go func() {
			inChan, _ := ipcrenderer.Endpoint{}.Listen(ChannelSignins)
			for {
				inMsg := <-inChan
				var update signinAnnouncement
				err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(&update)
				if err != nil {
					console.Log(err)
					continue
				}
				if !update.IsSignedIn {
					ClearToken()
				} else {
					SetToken(update.Token)
				}
			}
		}()
	}
}
