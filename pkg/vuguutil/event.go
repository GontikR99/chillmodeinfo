// +build wasm, web

package vuguutil

import (
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/js"
	js2 "syscall/js"
)

type fauxVuguEvent struct {
	event js.Value
	env vugu.EventEnv
}

func NewVuguEvent(event js2.Value, env vugu.EventEnv) vugu.DOMEvent {
	return &fauxVuguEvent{event: js.Value(event), env: env}
}

func (f *fauxVuguEvent) Prop(keys ...string) interface{} {
	cv := f.event
	for _, k := range keys {
		cv=cv.Get(k)
		if cv.IsNull() || cv.IsUndefined() {
			return nil
		}
	}
	return cv
}

func (f *fauxVuguEvent) PropString(keys ...string) string {
	sv := f.Prop(keys...)
	if sv==nil {
		return ""
	} else {
		return sv.(js.Value).String()
	}
}

func (f *fauxVuguEvent) PropFloat64(keys ...string) float64 {
	sv := f.Prop(keys...)
	if sv==nil {
		return 0
	} else {
		return sv.(js.Value).Float()
	}
}

func (f *fauxVuguEvent) PropBool(keys ...string) bool {
	sv := f.Prop(keys...)
	if sv==nil {
		return false
	} else {
		return sv.(js.Value).Bool()
	}
}

func (f *fauxVuguEvent) EventSummary() map[string]interface{} {
	return nil
}

func (f *fauxVuguEvent) JSEvent() js.Value {
	return f.event
}

func (f *fauxVuguEvent) JSEventTarget() js.Value {
	return f.event.Get("target")
}

func (f *fauxVuguEvent) JSEventCurrentTarget() js.Value {
	return f.event.Get("currentTarget")
}

func (f *fauxVuguEvent) EventEnv() vugu.EventEnv {
	return f.env
}

func (f *fauxVuguEvent) PreventDefault() {
	f.event.Call("preventDefault")
}

func (f *fauxVuguEvent) StopPropagation() {
	f.event.Call("stopPropagation")
}

