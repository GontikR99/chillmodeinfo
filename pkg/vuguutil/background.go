// +build wasm, web

package vuguutil

import (
	"context"
	"github.com/vugu/vugu"
	"time"
)

type BackgroundLooper interface {
	vugu.Builder
	RunInBackground()
}

type BackgroundComponent struct {
	Ctx context.Context
	env vugu.EventEnv
	renderChan <-chan struct{}
	renderCallbackHandle CallbackHandle
	cancelFunc context.CancelFunc
}

func (c *BackgroundComponent) InitBackground(vCtx vugu.InitCtx, bg BackgroundLooper) {
	c.Ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.env = vCtx.EventEnv()
	rChan := make(chan struct{})
	c.renderCallbackHandle=OnRender(func() {
		rChan<-struct{}{}
	})
	c.renderChan=rChan
	go func() {
		bg.RunInBackground()
	}()
}

func (c *BackgroundComponent) Rendered() <-chan struct{} {
	return c.renderChan
}

func (c *BackgroundComponent) Deadline() (deadline time.Time, ok bool) {
	return c.Ctx.Deadline()
}

func (c *BackgroundComponent) Done() <-chan struct{} {
	return c.Ctx.Done()
}

func (c *BackgroundComponent) Err() error {
	return c.Ctx.Err()
}

func (c *BackgroundComponent) Value(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *BackgroundComponent) Destroy(vCtx vugu.DestroyCtx) {
	c.renderCallbackHandle.Release()
	c.cancelFunc()
}

func (c *BackgroundComponent) Env() vugu.EventEnv {
	return c.env
}