// +build wasm,web

package main

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/vugu/vugu"
)

type GuildDump struct {
	Owner *Root
	Update *update.UpdateEntry
	Error string
}

func (c *GuildDump) upload(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	go func() {
		var m []record.Member
		for _, v := range c.Update.Members {
			m=append(m, v)
		}
		_, err := restidl.Members.MergeMembers(context.Background(), m)
		if err!=nil {
			c.Error=err.Error()
			return
		}
		c.Owner.updateMembers(event.EventEnv())
		c.Owner.removeFromQueue(event.EventEnv(), c.Update)
	}()
}

func (c *GuildDump) dismiss(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.Owner.removeFromQueue(event.EventEnv(), c.Update)
}