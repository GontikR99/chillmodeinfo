// +build wasm,web

package eqspec

import (
	"bytes"
	"encoding/gob"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
)

func ListenForLogs() (<-chan []*LogEntry, func()) {
	outChan := make(chan []*LogEntry)
	inChan, done := ipcrenderer.Endpoint{}.Listen(ChannelEQLog)
	go func() {
		for {
			inMsg := <-inChan
			if inMsg == nil {
				close(outChan)
				return
			}
			var entries []*LogEntry
			err := gob.NewDecoder(bytes.NewReader(inMsg.Content())).Decode(&entries)
			if err != nil {
				console.Log(err)
			} else {
				outChan <- entries
			}
		}
	}()
	return outChan, done
}
