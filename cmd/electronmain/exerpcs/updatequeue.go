// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
)

var nextHandler rpcidl.UpdateQueueHandler

type updateQueueServerHandler struct{}

func (u updateQueueServerHandler) Poll() (map[int]*overlay.UpdateEntry, error) {
	return nextHandler.Poll()
}

func (u updateQueueServerHandler) Enqueue(m map[int]*overlay.UpdateEntry) error {
	return nextHandler.Enqueue(m)
}

func SetUpdateQueueHandler(next rpcidl.UpdateQueueHandler) {
	nextHandler=next
}

func init() {
	register(rpcidl.HandleUpdateQueue(updateQueueServerHandler{}))
}
