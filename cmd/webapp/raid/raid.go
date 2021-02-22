// +build wasm,web

package raid

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"sort"
	"strconv"
	"time"
)

type Raid struct {
	Raids []record.Raid
	raidOpen map[uint64]struct{}
	ctx context.Context
	ctxDone context.CancelFunc
}

type raidTableEntry struct {
	idx int
	raid record.Raid
	mainLine bool
	credited []string
	uncredited []string
}

func (c *Raid) generateTableEntries() []raidTableEntry {
	var entries []raidTableEntry
	for idx, raid := range c.Raids {
		entries = append(entries, raidTableEntry{idx, raid, true, nil, nil})
		if !c.isRaidCollapsed(raid.GetRaidId()) {
			creditMap := make(map[string]struct{})
			for _, v := range raid.GetCredited() {
				creditMap[v]= struct{}{}
			}
			var credited []string
			var uncredited []string
			for _, v := range raid.GetAttendees() {
				if _, ok := creditMap[v]; ok {
					credited = append(credited, v)
				} else {
					uncredited=append(uncredited, v)
				}
			}
			sort.Sort(byValue(credited))
			sort.Sort(byValue(uncredited))
			entries=append(entries, raidTableEntry{idx, raid, false, credited, uncredited})
		}
	}
	return entries
}

func (c *Raid) isRaidCollapsed(raidId uint64) bool {
	_, ok := c.raidOpen[raidId]
	return !ok
}

func (c *Raid) toggleCollapsed(event vugu.DOMEvent, raidId uint64) {
	event.StopPropagation()
	event.PreventDefault()
	if _, ok := c.raidOpen[raidId]; ok {
		delete(c.raidOpen, raidId)
	} else {
		c.raidOpen[raidId]=struct{}{}
	}
}

func (c *Raid) updateDescription(submit ui.SubmitEvent, raid record.Raid) {
	go func() {
		if submit.Value()=="" {
			submit.Reject(errors.New("Give this raid a name!"))
			return
		}
		newRaid := record.NewBasicRaid(raid)
		newRaid.Description=submit.Value()
		_, err := restidl.Raid.Update(c.ctx, newRaid)
		if err!=nil {
			submit.Reject(err)
			return
		} else {
			submit.Accept(submit.Value())
			c.refreshRaids(submit.EventEnv())
		}
	}()
}

func rowStyle(idx int) string {
	if idx%2==0 {
		return "background-color: rgba(0,0,0,0.05)"
	} else {
		return ""
	}
}

func (c *Raid) recalculateRaid(event vugu.DOMEvent, raid record.Raid) {
	event.PreventDefault()
	event.StopPropagation()
	go func() {
		_, err := restidl.Raid.Update(c.ctx, raid)
		if err!=nil {
			toast.Error("Raids", err)
			return
		}
		c.refreshRaids(event.EventEnv())
	}()
}

func (c *Raid) updateDKP(submit ui.SubmitEvent, raid record.Raid) {
	go func() {
		dkpValue, err := strconv.ParseFloat(submit.Value(), 64)
		if err!=nil {
			submit.Reject(errors.New("Sorry, that's not a number"))
			return
		}
		newRaid := record.NewBasicRaid(raid)
		newRaid.DKPValue=dkpValue
		_, err = restidl.Raid.Update(c.ctx, newRaid)
		if err!=nil {
			submit.Reject(err)
			return
		} else {
			submit.Accept(submit.Value())
			c.refreshRaids(submit.EventEnv())
		}
	}()
}


func (c *Raid) refreshRaids(env vugu.EventEnv) {
	raids, err := restidl.Raid.Fetch(c.ctx)
	if err!=nil {
		toast.Error("raids", err)
		return
	}
	env.Lock()
	c.Raids=raids
	env.UnlockRender()
}

func (c *Raid) deleteRaid(event vugu.DOMEvent, raid record.Raid) {
	event.PreventDefault()
	event.StopPropagation()
	go func() {
		err := restidl.Raid.Delete(c.ctx, raid.GetRaidId())
		if err!=nil {
			toast.Error("raids", err)
			return
		}
		c.refreshRaids(event.EventEnv())
	}()
}

func (c *Raid) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
	c.raidOpen=make(map[uint64]struct{})
	go func() {
		for {
			c.refreshRaids(vCtx.EventEnv())
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(60*time.Second):
			}
		}
	}()
}

func (c *Raid) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}
