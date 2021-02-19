// +build wasm,web

package member

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"strings"
	"time"
)

type Member struct {
	Member record.Member
	LogEntries []record.DKPChangeEntry
	ctx context.Context
	ctxDone context.CancelFunc
}


func (c *Member) Init(vCtx vugu.InitCtx) {
	placeParts := strings.Split(place.GetPlace(), ":")
	if len(placeParts)>=1 {
		c.Member=&record.BasicMember{
			Name: placeParts[1],
		}
	} else {
		c.Member=&record.BasicMember{}
	}
	c.ctx, c.ctxDone = context.WithCancel(context.Background())
	go func() {
		mRec, err := restidl.Members.GetMember(c.ctx, c.Member.GetName())
		if err!=nil {
			toast.Error("member page", err)
		} else if mRec!=nil {
			vCtx.EventEnv().Lock()
			c.Member=mRec
			vCtx.EventEnv().UnlockRender()
		}
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(60*time.Second):
		}
	}()
	go func() {
		for {
			entries, err := restidl.DKPLog.Retrieve(c.ctx, c.Member.GetName())
			if err!=nil {
				toast.Error("member page", err)
			} else {
				vCtx.EventEnv().Lock()
				c.LogEntries=entries
				vCtx.EventEnv().UnlockRender()
			}
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(60*time.Second):
			}
		}
	}()
}

func (c *Member) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}