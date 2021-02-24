// +build wasm,web

package dom

import "syscall/js"

type Element interface {
	js.Wrapper
	AddEventListener(name string, jsFunc js.Func)

	Focus()
	Blur()

	AppendChild(Element)
	Remove()

	Get(key string) js.Value
	Set(key string, value interface{})

	SetAttribute(key string, value interface{})
	GetAttribute(key string) js.Value
	RemoveAttribute(key string)
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

func (j *jsDOMElement) SetAttribute(key string, value interface{}) {
	j.jsValue.Call("setAttribute", key, value)
}

func (j *jsDOMElement) GetAttribute(key string) js.Value {
	return j.jsValue.Call("getAttribute", key)
}

func (j *jsDOMElement) RemoveAttribute(key string) {
	j.jsValue.Call("removeAttribute", key)
}

func (j *jsDOMElement) Focus() {
	j.jsValue.Call("focus")
}

func (j *jsDOMElement) Blur() {
	j.jsValue.Call("blur")
}

func (j *jsDOMElement) Get(key string) js.Value {
	return j.jsValue.Get(key)
}

func (j *jsDOMElement) Set(key string, value interface{}) {
	j.jsValue.Set(key, value)
}

func WrapElement(value js.Value) Element {
	return &jsDOMElement{value}
}