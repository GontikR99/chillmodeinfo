// +build wasm,web

package leaderboard

import (
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"sort"
	"time"
)

type ClassCard struct {
	Class string
	Board *Leaderboard
}

func (c *ClassCard) filteredMembers() []record.Member {
	var res []record.Member
	for _, v := range c.Board.Members {
		if c.Board.ShouldShow(v) && c.Class == v.GetClass() {
			res=append(res, v)
		}
	}
	sort.Sort(byDKPthenName(res))
	return res
}

type byDKPthenName []record.Member

func (b byDKPthenName) Len() int {return len(b)}
func (b byDKPthenName) Swap(i, j int) {b[i], b[j] = b[j], b[i]}
func (b byDKPthenName) Less(i, j int) bool {
	if b[i].GetDKP()>b[j].GetDKP() {
		return true
	} else {
		return b[i].GetName() < b[j].GetName()
	}
}

func (c *ClassCard) updateDKP(event ui.SubmitEvent, member record.Member) {
	go func() {
		<-time.After(1*time.Second)
		event.Accept(event.Value())
	}()
}