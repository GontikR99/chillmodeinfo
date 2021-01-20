// +build wasm

package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/web/app/admin"
	"github.com/GontikR99/chillmodeinfo/web/app/home"
	"github.com/GontikR99/chillmodeinfo/web/app/leaderboard"
	"github.com/GontikR99/chillmodeinfo/web/app/login"
	"github.com/vugu/vugu"
)

type Root struct {
	LastPlace string
	Body vugu.Builder
}

type routeEntry struct {
	Place string
	DisplayName string
	Icon string
	BodyGen func()vugu.Builder
}

var routes =[]routeEntry{
	{"", "", "", func()vugu.Builder{return &home.Home{}}},
	{"login", "", "", func()vugu.Builder{return &login.Login{}}},
	{"register", "", "", func()vugu.Builder{return &login.Register{}}},
	{"leaderboard", "Leaderboard", "target", func()vugu.Builder{return &leaderboard.Leaderboard{}}},
	{"admin", "Admin", "terminal", func()vugu.Builder{return &admin.Admin{}}},
}

func (c *Root) Init(ctx vugu.InitCtx) {
	c.Body = &home.Home{}
}

func (c *Root) Compute(ctx vugu.ComputeCtx) {
	curPlace := place.GetPlace()
	if curPlace==c.LastPlace {
		return
	}
	for _, route := range routes {
		if route.Place==curPlace {
			c.Body = route.BodyGen()
			c.LastPlace = curPlace
			return
		}
	}
	place.NavigateTo(ctx.EventEnv(), "")
}