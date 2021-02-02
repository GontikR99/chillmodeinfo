// +build wasm

package ipcrenderer

import (
	"github.com/GontikR99/chillmodeinfo/internal/rpc"
	"syscall/js"
)

var ipcRenderer = js.Global().Get("ipcRenderer")

func Listen(channelName string) <-chan rpc.Message {
	resultChan := make(chan rpc.Message)
	if !ipcRenderer.IsUndefined() {
		ipcRenderer.Call("on", rpc.Prefix+channelName, js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			event := args[0]
			data := args[1].String()
			resultChan <- &electronMessage{
				event:   event,
				content: []byte(data),
			}
			return nil
		}))
	}
	return resultChan
}

func Send(channelName string, content []byte) {
	if !ipcRenderer.IsUndefined() {
		ipcRenderer.Call("send", rpc.Prefix+channelName, string(content))
	}
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
