// +build wasm,web

package login

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
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
		}
	}()
}

func (c *Standalone) Destroy(ctx vugu.DestroyCtx) {
	close(c.workingChan)
}
