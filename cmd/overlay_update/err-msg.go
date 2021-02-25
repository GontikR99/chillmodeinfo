// +build wasm,web

package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/vugu/vugu"
)

type ErrMsg struct {
	Owner  *Root
	Update *update.UpdateEntry
}

func (c *ErrMsg) dismiss(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.Owner.removeFromQueue(event.EventEnv(), c.Update)
}
