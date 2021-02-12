// +build wasm,web

package admin

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/cmd/webapp/localprofile"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/vugu/vugu"
	"strings"
	"time"
)

type Admin struct {
	DisplayName *FormTextInput
}

func (c *Admin) Init(ctx vugu.InitCtx) {
	c.DisplayName=&FormTextInput{
		Env: ctx.EventEnv(),
	}
	self := localprofile.GetProfile()
	if self!=nil {
		c.DisplayName.Value=self.GetDisplayName()
	}
}

func (c *Admin) submittable() bool {
	return c.DisplayName.Value!=""
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
	Env vugu.EventEnv
}

func (f *FormTextInput) StringValue() string {
	return f.Value
}

func (f *FormTextInput) SetStringValue(s string) {
	f.Value=strings.TrimSpace(s)
	go func() {
		f.Env.Lock()
		f.Env.UnlockRender()
	}()
}

