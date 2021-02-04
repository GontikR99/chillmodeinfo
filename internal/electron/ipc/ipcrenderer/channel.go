// +build wasm

package ipcrenderer

import (
	"encoding/hex"
	"github.com/GontikR99/chillmodeinfo/internal/msgcomm"
	"net/rpc"
	"syscall/js"
)

var ipcRenderer = js.Global().Get("ipcRenderer")
type Endpoint struct{}

func (i Endpoint) Listen(channelName string) (<-chan msgcomm.Message, func()) {
	resultChan := make(chan msgcomm.Message)
	recvFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		event := args[0]
		data , _ := hex.DecodeString(args[1].String())
		resultChan <- &electronMessage{
			event:   event,
			content: data,
		}
		return nil
	})
	ipcRenderer.Call("on", msgcomm.Prefix+channelName, recvFunc)
	return resultChan, func() {
		ipcRenderer.Call("removeListener", msgcomm.Prefix+channelName, recvFunc)
		recvFunc.Release()
		close(resultChan)
	}
}

func (i Endpoint) Send(channelName string, content []byte) {
	ipcRenderer.Call("send", msgcomm.Prefix+channelName, hex.EncodeToString(content))
}

type electronMessage struct {
	event js.Value
	content []byte
}

func (e *electronMessage) Content() []byte {
	return e.content
}

func (e *electronMessage) Sender() string {
	return "mainProcess"
}

func (e *electronMessage) Reply(channelName string, data []byte) {
	e.event.Call("reply", msgcomm.Prefix+channelName, hex.EncodeToString(data))
}

// Renderer side RPC client
var Client *rpc.Client

func init() {
	if !ipcRenderer.IsUndefined() {
		Client = msgcomm.NewClient(msgcomm.RpcMainChannel, &Endpoint{})
	}
}
