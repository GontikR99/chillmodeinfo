// +build wasm,web

package member

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/modal"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"strconv"
	"strings"
	"time"
)

type Member struct {
	Member     record.Member
	LogEntries []record.DKPChangeEntry
	ctx        context.Context
	ctxDone    context.CancelFunc
}

func descriptionStyle(entry record.DKPChangeEntry) string {
	if eqspec.IsItem(entry.GetDescription()) {
		return "font-weight: bold;"
	} else {
		return ""
	}
}

func (c *Member) updateDescription(submit ui.SubmitEvent, currentEntry record.DKPChangeEntry) {
	newEntry := record.NewBasicDKPChangeEntry(currentEntry)
	newEntry.Description = submit.Value()
	c.processLogUpdate(submit, newEntry)
}

func (c *Member) updateDelta(submit ui.SubmitEvent, currentEntry record.DKPChangeEntry) {
	newEntry := record.NewBasicDKPChangeEntry(currentEntry)
	deltaValue, err := strconv.ParseFloat(submit.Value(), 64)
	if err != nil {
		submit.Reject(errors.New("Please input a number"))
		return
	}
	newEntry.Delta = deltaValue
	c.processLogUpdate(submit, newEntry)
}

func (c *Member) processLogUpdate(submit ui.SubmitEvent, newEntry record.DKPChangeEntry) {
	go func() {
		update, err := restidl.DKPLog.Update(c.ctx, newEntry)
		if err != nil {
			submit.Reject(err)
			return
		} else {
			for idx, oldEntry := range c.LogEntries {
				if oldEntry.GetEntryId() == newEntry.GetEntryId() {
					submit.EventEnv().Lock()
					c.LogEntries[idx] = update
					submit.EventEnv().UnlockRender()
					break
				}
			}
			go c.reloadMember(submit.EventEnv())
			go c.reloadLogs(submit.EventEnv())
			submit.Accept(submit.Value())
		}
	}()
}

func (c *Member) cancelEntry(event vugu.DOMEvent, entry record.DKPChangeEntry) {
	event.StopPropagation()
	event.PreventDefault()
	go func() {
		if !modal.Verify("member", "Remove: "+entry.GetDescription(), "Are you sure you wish to remove this DKP change?", "Remove") {
			return
		}
		err := restidl.DKPLog.Remove(c.ctx, entry.GetEntryId())
		if err != nil {
			toast.Error("member page", err)
		} else {
			go c.reloadMember(event.EventEnv())
			go c.reloadLogs(event.EventEnv())
		}
	}()
}

func (c *Member) reloadMember(env vugu.EventEnv) {
	mRec, err := restidl.Members.GetMember(c.ctx, c.Member.GetName())
	if err != nil {
		toast.Error("member page", err)
	} else if mRec != nil {
		env.Lock()
		c.Member = mRec
		env.UnlockRender()
	}
}

func (c *Member) reloadLogs(env vugu.EventEnv) {
	entries, err := restidl.DKPLog.Retrieve(c.ctx, c.Member.GetName())
	if err != nil {
		toast.Error("member page", err)
	} else {
		env.Lock()
		c.LogEntries = entries
		env.UnlockRender()
	}
}

func (c *Member) Init(vCtx vugu.InitCtx) {
	placeParts := strings.Split(place.GetPlace(), ":")
	if len(placeParts) >= 1 {
		c.Member = &record.BasicMember{
			Name: placeParts[1],
		}
	} else {
		c.Member = &record.BasicMember{}
	}
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
	go func() {
		for {
			c.reloadMember(vCtx.EventEnv())
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(60 * time.Second):
			}
		}
	}()
	go func() {
		for {
			c.reloadLogs(vCtx.EventEnv())
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(60 * time.Second):
			}
		}
	}()
}

func (c *Member) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}
