// +build wasm,web

package roster

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"sort"
)

type Roster struct {
	membership []record.Member
	hideAlts bool
	sortOrder []sortOrder
}

func (c *Roster) reloadMembers(env vugu.EventEnv) {
	members, err := restidl.Members.GetMembers(context.Background())
	if err!=nil {
		toast.Error("members", err)
	}
	c.membership = []record.Member{}
	for _, v := range members {
		c.membership = append(c.membership, v)
	}
	c.resortMembers(env)
}

func (c *Roster) resortMembers(env vugu.EventEnv) {
	env.Lock()
	defer env.UnlockRender()
	for i:=len(c.sortOrder)-1;i>=0;i-- {
		ordering := c.sortOrder[i]
		var s sort.Interface
		switch ordering.index {
		case sortByName:
			s=byName(c.membership)
		case sortByClass:
			s=byClass(c.membership)
		case sortByLevel:
			s=byLevel(c.membership)
		case sortByAlt:
			s=byAlt(c.membership)
		case sortByOwner:
			s=byOwner(c.membership)
		case sortByActive:
			s=byActive(c.membership)
		}
		switch ordering.direction {
		case sortAscending:
			// default
		case sortDescending:
			s=reverse(s)
		}
		sort.Stable(s)
	}
}

func (c *Roster) updateSort(env vugu.DOMEvent, index sortIndex) {
	env.StopPropagation()
	env.PreventDefault()

	if c.sortOrder[0].index==index {
		c.sortOrder[0].direction=sortDirection(1-c.sortOrder[0].direction)
	} else {
		newSortOrder := []sortOrder{{index, sortAscending}}
		for _, v := range c.sortOrder {
			if v.index == index {
				continue
			} else {
				newSortOrder = append(newSortOrder, v)
			}
		}
		c.sortOrder = newSortOrder
	}
	console.Log(c.sortOrder)
	go c.resortMembers(env.EventEnv())
}

func (c *Roster) Init(ctx vugu.InitCtx) {
	c.sortOrder=[]sortOrder{{sortByName, sortAscending}}
	go c.reloadMembers(ctx.EventEnv())
}

type sortIndex int
type sortDirection int
const (
	sortByName=sortIndex(iota)
	sortByLevel
	sortByClass
	sortByAlt
	sortByOwner
	sortByActive

	sortAscending=sortDirection(0)
	sortDescending=sortDirection(1)
)

type sortOrder struct {
	index sortIndex
	direction sortDirection
}


type byName []record.Member
func (b byName) Len() int {return len(b)}
func (b byName) Less(i, j int) bool {return b[i].GetName() < b[j].GetName()}
func (b byName) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

type byLevel []record.Member
func (b byLevel) Len() int {return len(b)}
func (b byLevel) Less(i, j int) bool {return b[i].GetLevel() < b[j].GetLevel()}
func (b byLevel) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

type byClass []record.Member
func (b byClass) Len() int {return len(b)}
func (b byClass) Less(i, j int) bool {return b[i].GetClass() < b[j].GetClass()}
func (b byClass) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

type byAlt []record.Member
func (b byAlt) Len() int {return len(b)}
func (b byAlt) Less(i, j int) bool {return !b[i].IsAlt() && b[j].IsAlt()}
func (b byAlt) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

type byOwner []record.Member
func (b byOwner) Len() int {return len(b)}
func (b byOwner) Less(i, j int) bool {return b[i].GetOwner() < b[j].GetOwner()}
func (b byOwner) Swap(i, j int) {b[i], b[j] = b[j], b[i]}

type byActive []record.Member
func (b byActive) Len() int {return len(b)}
func (b byActive) Less(i,j int) bool {return b[i].GetLastActive().Before(b[j].GetLastActive())}
func (b byActive) Swap(i,j int) {b[i], b[j] = b[j], b[i]}

func reverse(s sort.Interface) sort.Interface {
	return &reversed{s}
}
type reversed struct {
	orig sort.Interface
}
func (r *reversed) Len() int {return r.orig.Len()}
func (r *reversed) Less(i,j int) bool {return r.orig.Less(j, i)}
func (r *reversed) Swap(i,j int) {r.orig.Swap(j,i)}