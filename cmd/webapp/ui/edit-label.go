// +build wasm,web

package ui

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"syscall/js"
	"time"
)

var inputIdGen=0

type EditLabel struct {
	InputId string
	Value   string

	Editable bool
	Submit SubmitHandler

	AttrMap vugu.AttrMap
	EditStyle string
	TextStyle string

	editState int
	editValue string

	ctx context.Context
	ctxDone context.CancelFunc

	ctxEditing context.Context
	ctxEditingDone context.CancelFunc
}

const (
	editStateDisplay=iota
	editStateEditing
	editStateSubmitting
)

func (c *EditLabel) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())

	inputIdGen++
	c.InputId=fmt.Sprintf("editlabel-%d", inputIdGen)

	c.editState=editStateDisplay
	if c.EditStyle=="" {
		c.EditStyle="width: 20rem;"
	}
}

func (c *EditLabel) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
}

func (c *EditLabel) adaptEditor(env vugu.EventEnv) {
	c.ctxEditing, c.ctxEditingDone = context.WithCancel(c.ctx)
	keyupFunc :=js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jEvent:=args[0]
		code:=jEvent.Get("code").String()
		switch code {
		case "Enter":
			jEvent.Call("preventDefault")
			jEvent.Call("stopPropagation")
			jEvent.Get("target").Call("blur")

			env.Lock()
			c.acceptEditRaw(env)
			env.UnlockRender()
		case "Escape":
			jEvent.Call("preventDefault")
			jEvent.Call("stopPropagation")
			jEvent.Get("target").Call("blur")

			env.Lock()
			c.cancelEditRaw(env)
			env.UnlockRender()
		default:
		}
		return nil
	})
	go func() {
		applied:=false
		for !applied {
			elem := document.GetElementById(c.InputId)
			if elem!=nil {
				elem.AddEventListener("keyup", keyupFunc)
				elem.JSValue().Call("select")
				applied=true
			} else {
				select {
				case <-c.ctxEditing.Done():
					applied=true
				case <-time.After(10*time.Millisecond):
				}
			}
		}
		<-c.ctxEditing.Done()
		keyupFunc.Release()
	}()
}

func (c *EditLabel) startEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateDisplay {
		c.editValue = c.Value
		c.editState = editStateEditing
		c.adaptEditor(event.EventEnv())
	}
}

func (c *EditLabel) editChange(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateEditing {
		c.editValue = event.JSEventTarget().Get("value").String()
	}
}

func (c *EditLabel) cancelEditRaw(env vugu.EventEnv) {
	if c.editState==editStateEditing {
		c.editState = editStateDisplay
		c.ctxEditingDone()
	}
}

func (c *EditLabel) cancelEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.cancelEditRaw(event.EventEnv())
}

func (c *EditLabel) acceptEditRaw(env vugu.EventEnv) {
	if c.editState==editStateEditing {
		c.editState = editStateSubmitting
		c.ctxEditingDone()

		if c.Submit!=nil {
			c.Submit.SubmitHandle(&completedEdit{
				label: c,
				env:   env,
			})
		} else {
			c.Value = c.editValue
			c.editState=editStateDisplay
		}
	}
}

func (c *EditLabel) acceptEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	c.acceptEditRaw(event.EventEnv())
}

type completedEdit struct {
	label *EditLabel
	env vugu.EventEnv
}

func (c *completedEdit) Value() string {return c.label.editValue}
func (c *completedEdit) EventEnv() vugu.EventEnv {return c.env}
func (c *completedEdit) Accept(value string) {
	go func() {
		c.env.Lock()
		if c.label.editState == editStateSubmitting {
			c.label.Value = value
			c.label.editState = editStateDisplay
		}
		c.env.UnlockRender()
	}()
}

func (c *completedEdit) Reject(err error) {
	toast.Error("edit", err)
	go func() {
		c.env.Lock()
		if c.label.editState == editStateSubmitting {
			c.label.editState = editStateEditing
		}
		c.env.UnlockRender()
		c.label.adaptEditor(c.env)
	}()
}

type SubmitEvent interface {
	Value() string
	EventEnv() vugu.EventEnv

	Accept(value string)
	Reject(error)
}

//vugugen:event Submit
