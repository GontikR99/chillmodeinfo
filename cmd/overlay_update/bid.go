// +build wasm,web

package main

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/vugu/vugu"
	"math"
	"sort"
	"strconv"
)

type Bid struct {
	Owner *Root
	Update *update.UpdateEntry
	Error string

	Bidder string
	ItemName string
	Bid float64
}

func (c *Bid) Init(vCtx vugu.InitCtx) {
	c.ItemName=c.Update.ItemName
	c.Bidder=c.Update.Bidder
	c.Bid=c.Update.Bid
}

func (c *Bid) submittable() bool {
	_, bidderPresent := c.Owner.membership[c.Bidder]
	return c.Bid!=0 && c.ItemName!="" && bidderPresent
}

func (c *Bid) submit(event vugu.DOMEvent) {
	go func() {
		err := restidl.DKPLog.Append(context.Background(), &record.BasicDKPChangeEntry{
			Target:      c.Bidder,
			Delta:       -math.Abs(c.Bid),
			Description: c.ItemName,
		})
		if err!=nil {
			event.EventEnv().Lock()
			c.Error=err.Error()
			event.EventEnv().UnlockRender()
			return
		} else {
			c.Owner.removeFromQueue(event.EventEnv(), c.Update)
			c.Owner.updateMembers(event.EventEnv())
		}
	}()
}

func (c *Bid) dismiss(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.Owner.removeFromQueue(event.EventEnv(), c.Update)
}

func (c *Bid) updateItem(event ui.ChangeEvent) {
	go func() {
		event.Env().Lock()
		c.ItemName=event.Value()
		event.Env().UnlockRender()
	}()
}

func (c *Bid) memberList() []string {
	members:=[]string{""}
	for _, v := range c.Owner.membership {
		members=append(members, v.GetName())
	}
	sort.Sort(byValueFold(members))
	return members
}

func (c *Bid) updateBidder(event ui.ChangeEvent) {
	go func() {
		event.Env().Lock()
		c.Bidder = event.Value()
		event.Env().UnlockRender()
	}()
}

func (c *Bid) updateDKP(event vugu.DOMEvent) {
	v, err := strconv.ParseFloat(event.PropString("target", "value"), 64)
	if err==nil {
		c.Bid=v
		event.JSEventTarget().Set("value", fmt.Sprintf("%.1f", c.Bid))
	} else {
		c.Error="Value must be a number"
		event.JSEventTarget().Call("select")
	}
}

func (c *Bid) suggest(event ui.SuggestionEvent) {
	event.Propose(eqspec.SuggestCompletions(event.Value()))
}