// +build wasm,web

package leaderboard

import (
	"context"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/ui"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/vugu/vugu"
	"sort"
	"strconv"
)

type ClassCard struct {
	Class string
	Board *Leaderboard
}

func (c *ClassCard) updateRecruitmentTarget(event ui.SubmitEvent) {
	go func() {
		newTarget, err := strconv.Atoi(event.Value())
		if err != nil {
			event.Reject(errors.New("Please enter a number"))
			return
		}
		if newTarget<0 {
			event.Reject(errors.New("Please enter a positive number"))
			return
		}
		newRecruitmentTarget := &record.BasicRecruitmentTarget{
			Class:  c.Class,
			Target: uint(newTarget),
		}
		err = restidl.Recruit.Update(c.Board.ctx, newRecruitmentTarget)
		if err!=nil {
			event.Reject(err)
			return
		}
		event.EventEnv().Lock()
		newTargetList:=[]record.RecruitmentTarget{newRecruitmentTarget}
		for _,v := range c.Board.RecruitmentTargets {
			if v.GetClass()!=newRecruitmentTarget.GetClass() {
				newTargetList = append(newTargetList, v)
			}
		}
		c.Board.RecruitmentTargets=newTargetList
		event.EventEnv().UnlockRender()
		event.Accept(event.Value())
	}()
}

func (c *ClassCard) activeCount() int {
	count:=0
	for _, member := range c.Board.Members {
		if record.IsActive(member) && member.GetClass()==c.Class && !member.IsAlt() {
			count++
		}
	}
	return count
}

func (c *ClassCard) desiredCount() int {
	for _, v := range c.Board.RecruitmentTargets {
		if v.GetClass()==c.Class {
			return int(v.GetTarget())
		}
	}
	return 0
}

func (c *ClassCard) filteredMembers() []record.Member {
	var res []record.Member
	for _, v := range c.Board.Members {
		if c.Board.ShouldShow(v) && c.Class == v.GetClass() {
			res = append(res, v)
		}
	}
	sort.Sort(byDKPthenName(res))
	return res
}

type byDKPthenName []record.Member

func (b byDKPthenName) Len() int      { return len(b) }
func (b byDKPthenName) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byDKPthenName) Less(i, j int) bool {
	if b[i].GetDKP() != b[j].GetDKP() {
		return b[i].GetDKP() > b[j].GetDKP()
	} else {
		return b[i].GetName() < b[j].GetName()
	}
}

func (c *ClassCard) jumpToMember(event vugu.DOMEvent, member record.Member) {
	event.PreventDefault()
	event.StopPropagation()
	place.NavigateTo(event.EventEnv(), "member:"+member.GetName())
}

func (c *ClassCard) updateDKP(event ui.SubmitEvent, member record.Member) {
	go func() {
		newValue, err := strconv.ParseFloat(event.Value(), 64)
		if err != nil {
			event.Reject(err)
			return
		}
		delta := &record.BasicDKPChangeEntry{
			Target:      member.GetName(),
			Delta:       newValue - member.GetDKP(),
			Description: "Leaderboard edit",
			RaidId:      0,
		}
		err = restidl.DKPLog.Append(context.Background(), delta)
		if err != nil {
			event.Reject(err)
			return
		}

		updatedMember, err := restidl.Members.GetMember(context.Background(), member.GetName())
		if err != nil {
			event.Reject(err)
			return
		}
		event.Accept(fmt.Sprintf("%.1f", updatedMember.GetDKP()))
		event.EventEnv().Lock()
		c.Board.Members[updatedMember.GetName()] = updatedMember
		event.EventEnv().UnlockRender()
	}()
}
