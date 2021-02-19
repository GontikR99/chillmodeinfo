// +build wasm,web

package raid

import (
	"context"
	"github.com/vugu/vugu"
)

type Raid struct {
	ctx context.Context
	ctxDone context.CancelFunc
}

func (c *Raid) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
}

func (c *Raid) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}

func (c *Raid) refreshRaids(event DumpPostedEvent) {

}