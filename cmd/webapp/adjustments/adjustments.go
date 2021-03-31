// +build wasm,web

package adjustments

import (
	"context"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/modal"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"math"
	"sort"
	"strconv"
	"time"
)

const pageLength=25

type Adjustments struct {
	Logs []record.DKPChangeEntry
	MaxDisplay int

	ctx     context.Context
	ctxDone context.CancelFunc

	pendingMember      string
	pendingDKP         float64
	pendingDescription string
	appending          bool

	members map[string]record.Member
}

func (c *Adjustments) KeyList() []string {
	names := []string{""}
	for _, member := range c.members {
		names = append(names, member.GetName())
	}
	sort.Sort(byValue(names))
	return names
}

type byValue []string

func (b byValue) Len() int           { return len(b) }
func (b byValue) Less(i, j int) bool { return b[i] < b[j] }
func (b byValue) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func (c *Adjustments) updatePendingMember(event ui.ChangeEvent) {
	c.pendingMember = event.Value()
}

func (c *Adjustments) updatePendingDKP(event vugu.DOMEvent) {
	newval, err := strconv.ParseFloat(event.JSEventTarget().Get("value").String(), 64)
	if err != nil {
		toast.Error("dkp", errors.New("You must enter a number"))
		event.JSEventTarget().Set("value", fmt.Sprintf("%.1f", c.pendingDKP))
		event.JSEventTarget().Call("select")
		return
	}
	c.pendingDKP = newval
	event.JSEventTarget().Set("value", fmt.Sprintf("%.1f", c.pendingDKP))
}

func (c *Adjustments) updatePendingDescription(change ui.ChangeEvent) {
	c.pendingDescription = change.Value()
}

func descriptionStyle(entry record.DKPChangeEntry) string {
	if eqspec.IsItem(entry.GetDescription()) {
		return "font-weight: bold;"
	} else {
		return ""
	}
}

func (c *Adjustments) addEntry(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	c.appending = true
	go func() {
		defer func() {
			event.EventEnv().Lock()
			c.appending = false
			event.EventEnv().UnlockRender()
		}()
		<-time.After(100 * time.Millisecond)
		if c.pendingMember == "" {
			toast.Error("adjustments", errors.New("You must select a member"))
			return
		}
		if c.pendingDescription == "" {
			toast.Error("adjustments", errors.New("You must provide a description"))
			return
		}
		err := restidl.DKPLog.Append(c.ctx, &record.BasicDKPChangeEntry{
			Target:      c.pendingMember,
			Delta:       c.pendingDKP,
			Description: c.pendingDescription,
		})
		if err != nil {
			toast.Error("adjustments", err)
		}
		c.pendingDKP = 0
		c.pendingDescription = ""
		c.pendingMember = ""
		c.reloadLog(event.EventEnv())
	}()
}

func (c *Adjustments) cancelEntry(event vugu.DOMEvent, entry record.DKPChangeEntry) {
	event.PreventDefault()
	event.StopPropagation()
	go func() {
		if !modal.Verify("adjustments", "Remove: "+entry.GetDescription(), "Are you sure you wish to remove this DKP change?", "Remove") {
			return
		}
		err := restidl.DKPLog.Remove(c.ctx, entry.GetEntryId())
		if err != nil {
			toast.Error("adjustments", err)
			return
		} else {
			c.reloadLog(event.EventEnv())
		}
	}()
}

func (c *Adjustments) updateDescription(submit ui.SubmitEvent, currentEntry record.DKPChangeEntry) {
	newEntry := record.NewBasicDKPChangeEntry(currentEntry)
	newEntry.Description = submit.Value()
	c.processLogUpdate(submit, newEntry)
}

func (c *Adjustments) updateDelta(submit ui.SubmitEvent, currentEntry record.DKPChangeEntry) {
	newEntry := record.NewBasicDKPChangeEntry(currentEntry)
	deltaValue, err := strconv.ParseFloat(submit.Value(), 64)
	if err != nil {
		submit.Reject(errors.New("Please input a number"))
		return
	}
	newEntry.Delta = deltaValue
	c.processLogUpdate(submit, newEntry)
}

func (c *Adjustments) processLogUpdate(submit ui.SubmitEvent, newEntry record.DKPChangeEntry) {
	go func() {
		update, err := restidl.DKPLog.Update(c.ctx, newEntry)
		if err != nil {
			submit.Reject(err)
			return
		} else {
			for idx, oldEntry := range c.Logs {
				if oldEntry.GetEntryId() == newEntry.GetEntryId() {
					submit.EventEnv().Lock()
					c.Logs[idx] = update
					submit.EventEnv().UnlockRender()
					break
				}
			}
			go c.reloadLog(submit.EventEnv())
			submit.Accept(submit.Value())
		}
	}()
}

func (c *Adjustments) reloadLog(env vugu.EventEnv) {
	logs, err := restidl.DKPLog.Retrieve(c.ctx, "")
	if err != nil {
		toast.Error("adjustments", err)
		return
	}
	env.Lock()
	c.Logs = logs
	env.UnlockRender()
}

func (c *Adjustments) showMore(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	if c.MaxDisplay!=math.MaxInt32 {
		c.MaxDisplay += pageLength
	}
}

func (c *Adjustments) showAll(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	c.MaxDisplay = math.MaxInt32
}

func (c *Adjustments) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
	c.MaxDisplay = pageLength
	go func() {
		members, err := restidl.Members.GetMembers(c.ctx)
		if err != nil {
			toast.Error("adjustments", err)
			return
		}
		vCtx.EventEnv().Lock()
		c.members = members
		vCtx.EventEnv().UnlockRender()
	}()
	go func() {
		done := false
		for !done {
			c.reloadLog(vCtx.EventEnv())
			select {
			case <-c.ctx.Done():
				done = true
				break
			case <-time.After(60 * time.Second):
			}
		}
	}()
}

func (c *Adjustments) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}

func (c *Adjustments) suggest(event ui.SuggestionEvent) {
	event.Propose(eqspec.SuggestCompletions(event.Value()))
}
