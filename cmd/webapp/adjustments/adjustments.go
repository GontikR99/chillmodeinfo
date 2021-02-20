// +build wasm,web

package adjustments

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"time"
)

type Adjustments struct {
	Logs []record.DKPChangeEntry

	ctx context.Context
	ctxDone context.CancelFunc
}

func (c *Adjustments) reloadLog(env vugu.EventEnv) {
	logs, err := restidl.DKPLog.Retrieve(c.ctx, "")
	if err!=nil {
		toast.Error("adjustments", err)
		return
	}
	env.Lock()
	c.Logs=logs
	env.UnlockRender()
}

func (c *Adjustments) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
	go func() {
		for {
			c.reloadLog(vCtx.EventEnv())
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(60*time.Second):
			}
		}
	}()
}

func (c *Adjustments) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}