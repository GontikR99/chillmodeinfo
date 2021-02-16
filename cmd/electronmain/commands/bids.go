// +build wasm,electron

package commands

import (
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/overlaymap"
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
)

type BufferedLogEntries struct {
	buffer    []*eqfiles.LogEntry
	bw        browserwindow.BrowserWindow
	retreived bool
}

func (b *BufferedLogEntries) FetchBufferedMessages() []*eqfiles.LogEntry {
	if b.retreived {
		return nil
	} else {
		b.retreived =true
		return b.buffer
	}
}

func (b *BufferedLogEntries) DispatchMessages(entries []*eqfiles.LogEntry) {
	if b.retreived {
		eqfiles.SendLogsTo(entries, b.bw)
	} else {
		b.buffer = append(b.buffer, entries...)
	}
}

var lastListener eqfiles.ListenerHandle

func OpenBids(initialBuffer []*eqfiles.LogEntry) {
	logBuffer := &BufferedLogEntries{
		buffer: initialBuffer,
	}
	listenerHandle := eqfiles.RegisterLogsListener(func(newEntries []*eqfiles.LogEntry) {
		logBuffer.DispatchMessages(newEntries)
	})
	lastListener = listenerHandle

	go func() {
		page := overlaymap.Lookup("bid").Page
		bw := overlay.Launch(page, false)
		logBuffer.bw = bw
		//bw.JSValue().Get("webContents").Call("openDevTools", map[string]interface{} {
		//	"mode":"detach",
		//	"activate":"false",
		//})
		bw.OnClosed(func() {
			listenerHandle.Release()
		})

		rpcServer := exerpcs.NewServer()
		rpcidl.HandleLogEntryBuffer(logBuffer)(rpcServer)
		bw.ServeRPC(rpcServer)
	}()
}

func EndBids() {
	lastListener.Release()
}

func ClearBids() {
	page := overlaymap.Lookup("bid").Page
	go overlay.Close(page)
}