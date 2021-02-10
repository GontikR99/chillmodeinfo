// +build wasm,web

package dom

import "syscall/js"

type Element interface {
	js.Wrapper
	AddEventListener(name string, jsFunc js.Func)
	AppendChild(Element)
	Remove()
	SetAttribute(key string, value string)
}

type jsDOMElement struct {
	jsValue js.Value
}

func (j *jsDOMElement) JSValue() js.Value {
	return j.jsValue
}

func (j *jsDOMElement) AddEventListener(name string, jsFunc js.Func) {
	j.jsValue.Call("addEventListener", name, jsFunc)
}

func (j *jsDOMElement) AppendChild(child Element) {
	if child!=nil {
		j.JSValue().Call("appendChild", child.JSValue())
	}
}

func (j *jsDOMElement) Remove() {
	j.jsValue.Call("remove")
}

func (j *jsDOMElement) SetAttribute(key string, value string) {
	j.jsValue.Call("setAttribute", key, value)
}

func WrapElement(value js.Value) Element {
	return &jsDOMElement{value}
}