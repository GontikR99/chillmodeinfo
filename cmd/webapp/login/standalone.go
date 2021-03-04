// +build wasm,web

package login

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/place"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/pkg/vuguutil"
	"github.com/vugu/vugu"
	"time"
)

const PlaceExternalLogin = "external-login"

type Standalone struct {
	workingChan chan struct{}
}

func (c *Standalone) Init(ctx vugu.InitCtx) {
	signIns := rpcidl.GetSignIn(ipcrenderer.Client)
	c.workingChan = make(chan struct{})
	go func() {
		signIns.SignIn()
		for {
			select {
			case <-c.workingChan:
				return
			case <-time.After(1 * time.Second):
				signIns.PollSignIn()
			}
			if signins.SignedIn() {
				place.NavigateTo(ctx.EventEnv(), "")
			}
		}
	}()
}

func (c *Standalone) Destroy(ctx vugu.DestroyCtx) {
	close(c.workingChan)
}

func loginLink() string {
	clientId, _, _ := rpcidl.LookupSetting(ipcrenderer.Client, settings.ClientId)
	return sitedef.SiteURL+"/associate.html?"+clientId
}

func clickEvent(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	node := vuguutil.GetElementByNodeId("login-link-text")
	node.JSValue().Call("select")
	document.ExecCommand("copy")
}