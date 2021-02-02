// +build wasm

package ipcmain

import (
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"github.com/GontikR99/chillmodeinfo/internal/rpc"
	"syscall/js"
)

var ipcMain=electron.Get().Get("ipcMain")

func Listen(channelName string) <-chan rpc.Message {
	resultChan := make(chan rpc.Message)
	ipcMain.Call("on", rpc.Prefix+channelName, js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		event := args[0]
		data := args[1].String()
		resultChan <- &electronMessage{
			event: event,
			content: []byte(data),
		}
		return nil
	}))
	return resultChan
}

type electronMessage struct {
	event js.Value
	content []byte
}

func (e *electronMessage) Content() []byte {
	return e.content
}

func (e *electronMessage) Reply(channelName string, data []byte) {
	e.event.Call("reply", rpc.Prefix+channelName, string(data))
}

