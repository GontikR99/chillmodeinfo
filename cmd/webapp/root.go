// +build wasm,web

package main

import (
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/admin"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/home"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/leaderboard"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/localprofile"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/login"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/settings"
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/vugu/vugu"
)

type Root struct {
	LastPlace string
	Body      vugu.Builder
}

type routeEntry struct {
	Place       string
	DisplayName string
	Icon        string
	ShowInNav   func() bool
	BodyGen     func() vugu.Builder
}

var neverShow = func() bool { return false }
var alwaysShow = func() bool { return true }

var routes = []routeEntry{
	{"", "Home", "home", neverShow, func() vugu.Builder { return &home.Home{} }},
	{"leaderboard", "Leaderboard", "target", alwaysShow, func() vugu.Builder { return &leaderboard.Leaderboard{} }},
	{"admin", "Admin", "terminal", func() bool {
		if localprofile.GetProfile()==nil {
			return false
		} else {
			switch localprofile.GetProfile().GetAdminState() {
			case profile.StateAdminUnrequested: return true
			case profile.StateAdminRequested: return true
			case profile.StateAdminApproved: return true
			default: return false
			}
		}
	}, func() vugu.Builder { return &admin.Admin{} }},
}

func init() {
	if electron.IsPresent() {
		routes = append(routes, routeEntry{"settings", "Settings", "settings", alwaysShow, func() vugu.Builder { return &settings.Settings{} }})
		routes = append(routes, routeEntry{login.PlaceExternalLogin, "", "", neverShow, func() vugu.Builder { return &login.Standalone{} }})
	}
}

func (c *Root) Init(ctx vugu.InitCtx) {
	c.Body = &home.Home{}
}

func (c *Root) Compute(ctx vugu.ComputeCtx) {
	curPlace := place.GetPlace()
	if curPlace == c.LastPlace {
		return
	}
	for _, route := range routes {
		if route.Place == curPlace {
			c.Body = route.BodyGen()
			c.LastPlace = curPlace
			return
		}
	}
	place.NavigateTo(ctx.EventEnv(), "")
}
