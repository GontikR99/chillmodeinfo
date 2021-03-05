// +build wasm,web

package home

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/GontikR99/chillmodeinfo/pkg/vuguutil"
	"github.com/vugu/vugu"
	"sort"
	"time"
)

type Home struct {
	vuguutil.BackgroundComponent
	members map[string]record.Member
	recruitTargets []record.RecruitmentTarget
}

func (c *Home) Init(vCtx vugu.InitCtx) {
	c.InitBackground(vCtx, c)
}

func (c *Home) refreshStats() {
	go func() {
		members, err := restidl.Members.GetMembers(c)
		if err!=nil {
			toast.Error("home", err)
			return
		}
		targets, err := restidl.Recruit.Fetch(c)
		if err!=nil {
			toast.Error("home", err)
			return
		}
		c.Env().Lock()
		c.members=members
		c.recruitTargets=targets
		c.Env().UnlockRender()
	}()
}

func (c *Home) RunInBackground() {
	c.refreshStats()
	for {
		select {
		case <-c.Done():
			return
		case <-c.Rendered():
		case <-time.After(60*time.Second):
			c.refreshStats()
		}
	}
}

type rtargetInfo struct {
	Class string
	Target int
}
type byClass []*rtargetInfo

func (b byClass) Len() int {return len(b)}
func (b byClass) Less(i, j int) bool {return b[i].Class<b[j].Class}
func (b byClass) Swap(i, j int) {b[i],b[j]=b[j],b[i]}

func (c *Home) recruitmentTargets() []*rtargetInfo {
	if c.members==nil || c.recruitTargets==nil {
		return nil
	}
	var targets []*rtargetInfo
	for class, _ := range eqspec.ClassMap {
		rtarget := &rtargetInfo{
			Class:  class,
			Target: 0,
		}
		for _, member := range c.members {
			if !member.IsAlt() && record.IsActive(member) && member.GetClass()==class {
				rtarget.Target--
			}
		}
		for _, tinfo := range c.recruitTargets {
			if tinfo.GetClass()==class {
				rtarget.Target+=int(tinfo.GetTarget())
				break
			}
		}
		targets = append(targets, rtarget)
	}
	sort.Sort(byClass(targets))
	return targets
}