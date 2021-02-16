// +build wasm,electron

package localprofile

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"time"
)

func StartElectronPoll() {
	signins.OnStateChange(func() {
		if !signins.SignedIn() {
			currentProfile = nil
		}
	})
	go func() {
		for {
			if signins.SignedIn() {
				fetchedProfile, err := restidl.GetProfile().FetchMine(context.Background())
				if err != nil {
					console.Log(err)
				} else {
					currentProfile = fetchedProfile
				}
				<-time.After(5 * time.Second)
			} else {
				currentProfile = nil
				<-time.After(10 * time.Millisecond)
			}
		}
	}()
	go func() {
		for {
			buffer := new(bytes.Buffer)
			err := gob.NewEncoder(buffer).Encode(profileMessage{profile.NewBasicProfile(currentProfile)})
			if err == nil {
				browserwindow.Broadcast(channelProfile, buffer.Bytes())
			} else {
				console.Log(err)
			}
			<-time.After(10*time.Millisecond)
		}
	}()
}