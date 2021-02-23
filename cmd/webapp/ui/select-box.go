// +build wasm,web

package ui

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/vugu/vugu"
	"syscall/js"
	"time"
)

var selectBoxIdGen int

type SelectBox struct {
	AttrMap vugu.AttrMap
	Value   string
	Change  ChangeHandler
	Options []string

	idStr string

	ctx        context.Context
	doneFunc   context.CancelFunc
	changeFunc js.Func
}

func (c *SelectBox) onChange(event js.Value, env vugu.EventEnv) {
	c.Value = event.Get("target").Get("value").String()
	if c.Change != nil {
		c.Change.ChangeHandle(&selectBoxChangeEvent{
			value: c.Value,
			env:   env,
			sb:    c,
		})
	}
}

type selectBoxChangeEvent struct {
	value string
	env   vugu.EventEnv
	sb    *SelectBox
}

func (s *selectBoxChangeEvent) Value() string      { return s.value }
func (s *selectBoxChangeEvent) Env() vugu.EventEnv { return s.env }
func (s *selectBoxChangeEvent) SetValue(s2 string) {
	go func() {
		s.env.Lock()
		s.sb.Value = s2
		s.env.UnlockRender()
	}()
}

func (c *SelectBox) Init(vCtx vugu.InitCtx) {
	c.ctx, c.doneFunc = context.WithCancel(context.Background())

	selectBoxIdGen++
	c.idStr = fmt.Sprintf("select-box-%d", selectBoxIdGen)
	c.changeFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.onChange(args[0], vCtx.EventEnv())
		return nil
	})
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
			elt := document.GetElementById(c.idStr)
			if elt != nil {
				elt.JSValue().Set("onchange", c.changeFunc)
			}
			for idx, value := range c.Options {
				elt := document.GetElementById(fmt.Sprintf("%s-option-%d", c.idStr, idx))
				if elt != nil {
					if c.Value == value {
						elt.SetAttribute("selected", "selected")
					} else {
						elt.RemoveAttribute("selected")
					}
				}
			}
		}
	}()
}

func (c *SelectBox) Destroy(vCtx vugu.DestroyCtx) {
	c.doneFunc()
	c.changeFunc.Release()
}
