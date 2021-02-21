// +build wasm,electron

package logactions

import (
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/overlaymap"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
)

type BufferedLogEntries struct {
	buffer    []*eqspec.LogEntry
	bw        browserwindow.BrowserWindow
	retreived bool
}

func (b *BufferedLogEntries) FetchBufferedMessages() []*eqspec.LogEntry {
	if b.retreived {
		return nil
	} else {
		b.retreived =true
		return b.buffer
	}
}

func (b *BufferedLogEntries) DispatchMessages(entries []*eqspec.LogEntry) {
	if b.retreived {
		eqspec.SendLogsTo(entries, b.bw)
	} else {
		b.buffer = append(b.buffer, entries...)
	}
}

var lastListener eqspec.ListenerHandle

func OpenBids(initialBuffer []*eqspec.LogEntry) {
	logBuffer := &BufferedLogEntries{
		buffer: initialBuffer,
	}
	listenerHandle := eqspec.RegisterLogsListener(func(newEntries []*eqspec.LogEntry) {
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