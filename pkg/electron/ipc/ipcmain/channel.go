// +build wasm,electron

package ipcmain

import (
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc"
	"github.com/GontikR99/chillmodeinfo/pkg/msgcomm"
	"strconv"
	"syscall/js"
)

var ipcMain= electron.JSValue().Get("ipcMain")

func Listen(channelName string) (<-chan msgcomm.Message, func()) {
	resultChan := make(chan msgcomm.Message)
	recvFunc := js.FuncOf(func(_ js.Value, args []js.Value)interface{} {
		event := args[0]
		data, _ := ipc.Decode(args[1])
		resultChan <- &electronMessage{
			event: event,
			content: []byte(data),
		}
		return nil
	})
	ipcMain.Call("on", msgcomm.Prefix+channelName, recvFunc)
	return resultChan, func() {
		ipcMain.Call("removeListener", msgcomm.Prefix+channelName, recvFunc)
		recvFunc.Release()
		close(resultChan)
	}
}

type electronMessage struct {
	event js.Value
	content []byte
}

func (e *electronMessage) Content() []byte {
	return e.content
}

func (e *electronMessage) Sender() string {
	return strconv.Itoa(e.event.Get("sender").Get("id").Int())
}

func (e *electronMessage) Reply(channelName string, data []byte) {
	e.event.Call("reply", msgcomm.Prefix+channelName, ipc.Encode(data))
}
