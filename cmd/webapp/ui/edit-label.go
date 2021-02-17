// +build wasm,web

package ui

import (
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
)

type EditLabel struct {
	Value string

	Editable bool
	Submit SubmitHandler

	AttrMap vugu.AttrMap
	EditWidth string

	editState int
	editValue string
}

const (
	editStateDisplay=iota
	editStateEditing
	editStateSubmitting
)

func (c *EditLabel) Init(vCtx vugu.InitCtx) {
	c.editState=editStateDisplay
	if c.EditWidth=="" {
		c.EditWidth="20rem"
	}
}

func (c *EditLabel) startEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateDisplay {
		c.editValue = c.Value
		c.editState=editStateEditing
	}
}

func (c *EditLabel) editChange(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateEditing {
		c.editValue = event.JSEventTarget().Get("value").String()
	}
}

func (c *EditLabel) cancelEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateEditing {
		c.editState = editStateDisplay
	}
}

func (c *EditLabel) acceptEdit(event vugu.DOMEvent) {
	event.PreventDefault()
	event.StopPropagation()
	if c.editState==editStateEditing {
		c.editState = editStateSubmitting
		if c.Submit!=nil {
			c.Submit.SubmitHandle(&completedEdit{
				label: c,
				env:   event.EventEnv(),
			})
		} else {
			c.Value = c.editValue
			c.editState=editStateDisplay
			console.Log(c.editValue, c.Value)
		}
	}
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
	}()
}

type SubmitEvent interface {
	Value() string
	EventEnv() vugu.EventEnv

	Accept(value string)
	Reject(error)
}

//vugugen:event Submit
