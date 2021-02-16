// +build wasm,web

package login

import (
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/vugu/vugu"
)

type Login struct {
	SignedIn bool
}

func visStyle(value bool) string {
	if value {
		return "visibility: visible;"
	} else {
		return "visibility: hidden; position:absolute; top:0; left:0;"
	}
}

func (c *Login) Init(ctx vugu.InitCtx) {
	c.SignedIn = signins.SignedIn()

	signins.OnStateChange(func() {
		ctx.EventEnv().Lock()
		c.SignedIn = signins.SignedIn()
		ctx.EventEnv().UnlockRender()

		// Return to the home page
		place.NavigateTo(ctx.EventEnv(), "")
	})
}

func (c *Login) Compute(ctx vugu.ComputeCtx) {
	c.SignedIn = signins.SignedIn()
}

func (c *Login) SignIn(event vugu.DOMEvent) {
	event.PreventDefault()
	if electron.IsPresent() {
		place.NavigateTo(event.EventEnv(), PlaceExternalLogin)
	} else {
		signins.SignIn()
	}
}

func (c *Login) SignOut(event vugu.DOMEvent) {
	event.PreventDefault()
	signins.SignOut()
}
