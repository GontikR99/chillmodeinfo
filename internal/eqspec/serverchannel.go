// +build wasm,electron

package eqspec

import (
	"bytes"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
)

func SendLogsTo(entries []*LogEntry, bw browserwindow.BrowserWindow) {
	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(entries)
	if err != nil {
		console.Log(err)
		return
	}
	bw.Send(ChannelEQLog, buffer.Bytes())
}
