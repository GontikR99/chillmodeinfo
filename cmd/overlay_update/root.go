// +build wasm,web

package main

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
	"sort"
	"time"
)

type Root struct {
	queue map[int]*update.UpdateEntry
	membership map[string]record.Member
}

func (c *Root) removeFromQueue(env vugu.EventEnv, item *update.UpdateEntry) {
	go func() {
		env.Lock()
		delete(c.queue, item.SeqId)
		env.UnlockRender()
	}()
}

func (c *Root) updateMembers(env vugu.EventEnv) {
	go func() {
		members, err := restidl.Members.GetMembers(context.Background())
		if err!=nil {
			env.Lock()
			c.membership=members
			env.UnlockRender()
		}
	}()
}

var updateQueue=rpcidl.UpdateQueue(ipcrenderer.Client)

func (c *Root) Init(vCtx vugu.InitCtx) {
	c.queue=make(map[int]*update.UpdateEntry)
	c.membership=make(map[string]record.Member)
	go func() {
		c.updateMembers(vCtx.EventEnv())
		<-time.After(60*time.Second)
	}()
	go func() {
		for {
			<-time.After(100*time.Millisecond)
			newEntries, err := updateQueue.Poll()
			if err!=nil {
				console.Log(err)
				continue
			}
			if len(newEntries)!=0 {
				vCtx.EventEnv().Lock()
				for k, v := range newEntries {
					c.queue[k] = v
					localprofile.SetProfileIfAbsent(v.Self)
				}
				vCtx.EventEnv().UnlockRender()
			}
			remaining := map[int]*update.UpdateEntry{}
			for k,v := range c.queue {
				remaining[k]=v.Duplicate()
			}
			updateQueue.Enqueue(remaining)
		}
	}()
}

func (c *Root) enumerateQueue() []*update.UpdateEntry {
	var res []*update.UpdateEntry
	for _, v := range c.queue {
		res=append(res, v)
	}
	sort.Sort(bySeqNumDesc(res))
	return res
}

type bySeqNumDesc []*update.UpdateEntry

func (b bySeqNumDesc) Len() int {return len(b)}
func (b bySeqNumDesc) Less(i, j int) bool {return b[i].SeqId>b[j].SeqId}
func (b bySeqNumDesc) Swap(i, j int) {b[i], b[j] = b[j],b[i]}
