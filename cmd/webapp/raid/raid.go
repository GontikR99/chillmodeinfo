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
	"strconv"
	"time"
)

type Raid struct {
	Raids []record.Raid

	ctx context.Context
	ctxDone context.CancelFunc
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
