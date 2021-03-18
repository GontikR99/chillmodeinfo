// +build wasm,web

package admin

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
	"github.com/GontikR99/chillmodeinfo/pkg/modal"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"strings"
	"time"
)

type Admin struct {
	DisplayName *FormTextInput

	SyncResultText string
	SyncPending bool
}

func (c *Admin) Init(ctx vugu.InitCtx) {
	c.DisplayName = &FormTextInput{
		Env: ctx.EventEnv(),
	}
	self := localprofile.GetProfile()
	if self != nil {
		c.DisplayName.Value = self.GetDisplayName()
	}
}

func (c *Admin) submittable() bool {
	return c.DisplayName.Value != ""
}

func (c *Admin) RequestAdmin(event vugu.DOMEvent) {
	event.PreventDefault()
	if c.submittable() {
		go restidl.GetProfile().RequestAdmin(context.Background(), c.DisplayName.Value)
	} else {
		toast.PopupWithDuration("Required field", "Please tell us your name!", 30*time.Second)
	}
}

type FormTextInput struct {
	Value string
	Env   vugu.EventEnv
}

func (f *FormTextInput) StringValue() string {
	return f.Value
}

func (f *FormTextInput) SetStringValue(s string) {
	f.Value = strings.TrimSpace(s)
	go func() {
		f.Env.Lock()
		f.Env.UnlockRender()
	}()
}

func (c *Admin) syncDKP(event vugu.DOMEvent) {
	c.SyncPending=true
	go func() {
		if !modal.Verify("dkpsync", "Sync from Gamerlaunch", "Are you sure you wish to sync with Gamerlaunch now?", "Sync") {
			event.EventEnv().Lock()
			c.SyncPending=false
			event.EventEnv().UnlockRender()
			return
		}
		msg, err := restidl.DKPLog.Sync(context.Background())
		event.EventEnv().Lock()
		if err!=nil {
			c.SyncResultText=err.Error()
		} else {
			c.SyncResultText=msg
		}
		c.SyncPending=false
		event.EventEnv().UnlockRender()
	}()
}

func (c *Admin) syncDKPAttrs() []vugu.VGAttribute {
	if c.SyncPending {
		return []vugu.VGAttribute{{"", "disabled", "disabled"}}
	} else {
		return nil
	}
}