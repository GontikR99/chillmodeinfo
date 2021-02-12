// +build wasm, web

package localprofile

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/signins"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/vugu/vugu"
	"time"
)

var currentProfile profile.Entry

func GetProfile() profile.Entry {
	return currentProfile
}

func StartPoll(env vugu.EventEnv) {
	go func() {
		for {
			<-time.After(1*time.Second)
			if signins.SignedIn() {
				fetchedProfile, err := restidl.GetProfile().FetchMine(context.Background())
				if err!=nil {
					toast.Error("accounts", errors.New("Couldn't read your profile: "+err.Error()))
					signins.SignOut()
					env.Lock()
					currentProfile=nil
					env.UnlockRender()
					continue
				}
				if !profile.Equal(GetProfile(), fetchedProfile) {
					env.Lock()
					currentProfile=fetchedProfile
					env.UnlockRender()
				}
			} else {
				if currentProfile!=nil {
					env.Lock()
					currentProfile=nil
					env.UnlockRender()
				}
			}
		}
	}()
}

