// +build wasm,electron

package signins

import (
	"bytes"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"time"
)

func broadcastLoginState() {
	buffer := new(bytes.Buffer)
	token := ""
	if SignedIn() {
		token = *currentToken
	}
	err := gob.NewEncoder(buffer).Encode(&signinAnnouncement{
		IsSignedIn: SignedIn(),
		Token:      token,
	})
	if err != nil {
		console.Log(err)
		return
	}
	browserwindow.Broadcast(ChannelSignins, buffer.Bytes())
}

func init() {
	OnStateChange(func() {
		broadcastLoginState()
	})
	go func() {
		for {
			broadcastLoginState()
			<-time.After(100 * time.Millisecond)
		}
	}()
}
