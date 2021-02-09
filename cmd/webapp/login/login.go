// +build wasm,web

package login

import (
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/signin"
	"github.com/vugu/vugu"
)

type Login struct {
	UserName string
	SignedIn bool
}

func init() {
	signin.PrepareForSignin(sitedef.GoogleSigninClientId)
}

func (c *Login) Init(ctx vugu.InitCtx) {
	signin.OnStateChange(func(info *signin.UserInfo) {
		ctx.EventEnv().Lock()
		if info == nil {
			c.UserName = ""
			c.SignedIn = false
		} else {
			c.UserName = info.Name
			c.SignedIn = true
		}
		ctx.EventEnv().UnlockRender()

		// Return to the home page
		place.NavigateTo(ctx.EventEnv(), "")
	})
}

func (c *Login) Compute(ctx vugu.ComputeCtx) {
	c.SignedIn = signin.SignedIn()
}

func (c *Login) Rendered(ctx vugu.RenderedCtx) {
	if !c.SignedIn {
		signin.RenderSignin("chillmode-signin2")
	}
}

func (c *Login) SignOut(event vugu.DOMEvent) {
	event.PreventDefault()
	signin.SignOut()
}
