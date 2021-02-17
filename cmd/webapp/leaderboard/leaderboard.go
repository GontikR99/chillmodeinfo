// +build wasm,web

package leaderboard

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"sort"
	"time"
)

type Leaderboard struct {
	Cards []*ClassCard
	Members map[string]record.Member

	hideInactive bool
	hideAlts bool

	ctx context.Context
	ctxDone func()
}

func (c *Leaderboard) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())

	c.hideInactive=true
	c.hideAlts=true

	for k, _ := range eqfiles.ClassMap {
		c.Cards=append(c.Cards, &ClassCard{
			Class:   k,
			Board: c,
		})
	}
	sort.Sort(byClass(c.Cards))

	c.Members=make(map[string]record.Member)

	go func() {
		for {
			members, err := restidl.Members.GetMembers(c.ctx)
			if err != nil {
				toast.Error("leaderboard", err)
				return
			}
			vCtx.EventEnv().Lock()
			c.Members = members
			vCtx.EventEnv().UnlockRender()
			select {
				case <-c.ctx.Done():
					return
				case <-time.After(60*time.Second):
			}
		}
	}()
}

func (c *Leaderboard) hideInactiveChanged(event vugu.DOMEvent) {
	c.hideInactive=event.JSEventTarget().Get("checked").Bool()
}

func (c *Leaderboard) hideAltsChanged(event vugu.DOMEvent) {
	c.hideAlts=event.JSEventTarget().Get("checked").Bool()
}


func (c *Leaderboard) ShouldShow(member record.Member)bool {
	return (!c.hideInactive || record.IsActive(member)) &&
		(!c.hideAlts || !member.IsAlt())
}

func (c *Leaderboard) Destroy(ctx vugu.DestroyCtx) {
	c.ctxDone()
}


type byClass []*ClassCard
func (b byClass) Len() int {return len(b)}

func (b byClass) Less(i, j int) bool {return b[i].Class<b[j].Class}
func (b byClass) Swap(i, j int) {b[i],b[j] = b[j], b[i]}
