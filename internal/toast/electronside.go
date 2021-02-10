// +build wasm,electron

package toast

import (
	"bytes"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"time"
)

func PopupWithDuration(title string, messageText string, delay time.Duration) {
	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(&toastMessage{
		Title:   title,
		Body:    messageText,
		Timeout: delay,
	})
	if err==nil {
		browserwindow.Broadcast(channelToast, buffer.Bytes())
	}
}