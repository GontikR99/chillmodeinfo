// +build wasm,web

package login

import (
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/signins"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/document"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/vugu/vugu"
	"syscall/js"
)

type Login struct {
	Attached bool
	SignedIn bool
	ClickHandler *js.Func
}

func init() {
	signins.PrepareForSignin(sitedef.GoogleSigninClientId)
}

func visStyle(value bool) string {
	if value {
		return "visibility: visible;"
	} else {
		return "visibility: hidden; position:absolute; top:0; left:0;"
	}
}

func (c *Login) Init(ctx vugu.InitCtx) {
	c.SignedIn=signins.SignedIn()
	console.Logf("Signed in: %v", c.SignedIn)

	signins.OnStateChange(func() {
		ctx.EventEnv().Lock()
		c.SignedIn = signins.SignedIn()
		console.Logf("Signed in: %v", c.SignedIn)
		ctx.EventEnv().UnlockRender()

		// Return to the home page
		place.NavigateTo(ctx.EventEnv(), "")
	})
}

func (c *Login) Compute(ctx vugu.ComputeCtx) {
	c.SignedIn = signins.SignedIn()
}

func (c *Login) Rendered(ctx vugu.RenderedCtx) {
	if !c.Attached {
		signinButton := document.GetElementById("signin-button")
		c.Attached=true
		if electron.IsPresent() {
			c.ClickHandler=new(js.Func)
			*c.ClickHandler=js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				args[0].Call("preventDefault")
				go place.NavigateTo(ctx.EventEnv(), PlaceExternalLogin)
				return nil
			})
			signinButton.AddEventListener("click", *c.ClickHandler)
		} else {
			signins.Attach(signinButton)
		}
	}
}

func (c *Login) Destroy(ctx vugu.DestroyCtx) {
	if c.ClickHandler!=nil {
		c.ClickHandler.Release()
	}
}

func (c *Login) SignOut(event vugu.DOMEvent) {
	event.PreventDefault()
	signins.SignOut()
}
