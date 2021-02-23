// +build wasm,web

package main

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/overlay/update"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/vugu/vugu"
	"sort"
	"strconv"
	"strings"
)

type RaidDump struct {
	Owner *Root
	Update *update.UpdateEntry
	Error string
	RaidName string
	RaidValue float64

	prevRaidNames []string
}

func (c *RaidDump) Init(vCtx vugu.InitCtx) {
	go func() {
		raids, err := restidl.Raid.Fetch(context.Background())
		if err!=nil {
			return
		}
		raidSet := make(map[string]string)
		for _, raid := range raids {
			raidSet[strings.ToUpper(raid.GetDescription())]=raid.GetDescription()
		}

		vCtx.EventEnv().Lock()
		for _, v := range raidSet {
			c.prevRaidNames=append(c.prevRaidNames, v)
		}
		sort.Sort(byValueFold(c.prevRaidNames))
		vCtx.EventEnv().UnlockRender()
	}()
}

type byValueFold []string
func (b byValueFold) Len() int {return len(b)}
func (b byValueFold) Less(i, j int) bool {return strings.ToUpper(b[i]) < strings.ToUpper(b[j])}
func (b byValueFold) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

func (c *RaidDump) suggest(event ui.SuggestionEvent) {
	console.Log(event.Value())
	start := sort.Search(len(c.prevRaidNames), func(i int) bool {
		return strings.ToUpper(event.Value()) <= strings.ToUpper(c.prevRaidNames[i])
	})
	end := sort.Search(len(c.prevRaidNames), func(i int) bool {
		return strings.ToUpper(event.Value()+"\uFFFF") < strings.ToUpper(c.prevRaidNames[i])
	})
	console.Log(start, end)
	go event.Propose(c.prevRaidNames[start:end])
}

func (c *RaidDump) updateDescription(event ui.ChangeEvent) {
	go func() {
		event.Env().Lock()
		c.RaidName=event.Value()
		event.Env().UnlockRender()
	}()
}

func (c *RaidDump) updateDKP(event vugu.DOMEvent) {
	v, err := strconv.ParseFloat(event.PropString("target", "value"), 64)
	if err==nil {
		c.RaidValue=v
		event.JSEventTarget().Set("value", fmt.Sprintf("%.1f", c.RaidValue))
	} else {
		c.Error="Value must be a number"
		event.JSEventTarget().Call("select")
	}
}

func (c *RaidDump) upload(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	go func() {
		err := restidl.Raid.Add(context.Background(), &record.BasicRaid{
			Description: c.RaidName,
			Attendees:   c.Update.Attendees,
			DKPValue:    c.RaidValue,
		})
		if err!=nil {
			event.EventEnv().Lock()
			c.Error=err.Error()
			event.EventEnv().UnlockRender()
		} else {
			c.Owner.removeFromQueue(event.EventEnv(), c.Update)
		}
	}()
}

func (c *RaidDump) dismiss(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.Owner.removeFromQueue(event.EventEnv(), c.Update)
}