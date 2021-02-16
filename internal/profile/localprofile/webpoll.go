// +build wasm, web

package localprofile

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"time"
)

func StartWebPoll(env vugu.EventEnv) {
	signins.OnStateChange(func() {
		if !signins.SignedIn() {
			env.Lock()
			currentProfile=nil
			env.UnlockRender()
		}
	})
	go func() {
		if electron.IsPresent() {
			inChan, _ := ipcrenderer.Endpoint{}.Listen(channelProfile)
			for {
				inMsg := <- inChan
				var inProfileMsg profileMessage
				err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(&inProfileMsg)
				if err!=nil {
					toast.Error("localprofile", err)
					continue
				}
				newProfile := inProfileMsg.Value
				if currentProfile==nil {
					if newProfile==nil {
						continue
					} else if signins.SignedIn() {
						env.Lock()
						currentProfile = newProfile
						env.UnlockRender()
					}
				} else {
					if newProfile==nil {
						env.Lock()
						currentProfile = nil
						env.UnlockRender()
					} else {
						if !profile.Equal(currentProfile, newProfile) {
							env.Lock()
							currentProfile = newProfile
							env.UnlockRender()
						}
					}
				}
			}
		} else {
			for {
				if signins.SignedIn() {
					fetchedProfile, err := restidl.GetProfile().FetchMine(context.Background())
					if err != nil {
						toast.Error("accounts", errors.New("Couldn't read your profile: "+err.Error()))
						signins.SignOut()
						env.Lock()
						currentProfile = nil
						env.UnlockRender()
					} else if !profile.Equal(GetProfile(), fetchedProfile) {
						if signins.SignedIn() {
							env.Lock()
							currentProfile = fetchedProfile
							env.UnlockRender()
						}
					}
					<-time.After(5 * time.Second)
				} else {
					if currentProfile != nil {
						env.Lock()
						currentProfile = nil
						env.UnlockRender()
					}
					<-time.After(10 * time.Millisecond)
				}
			}
		}
	}()
}
