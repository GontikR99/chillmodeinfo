// +build wasm,web

package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
	"time"
)

type Root struct {
	queue map[int]*overlay.UpdateEntry
}

var updateQueue=rpcidl.UpdateQueue(ipcrenderer.Client)

func (c *Root) Init(vCtx vugu.InitCtx) {
	c.queue=make(map[int]*overlay.UpdateEntry)
	go func() {
		for {
			<-time.After(10*time.Millisecond)
			newEntries, err := updateQueue.Poll()
			if err!=nil {
				console.Log(err)
				continue
			}
			vCtx.EventEnv().Lock()
			for k, v := range newEntries {
				c.queue[k]=v
				console.Log(v)
			}
			vCtx.EventEnv().UnlockRender()
			remaining := map[int]*overlay.UpdateEntry {}
			for k,v := range c.queue {
				remaining[k]=v.Duplicate()
			}
			updateQueue.Enqueue(remaining)
		}
	}()
}