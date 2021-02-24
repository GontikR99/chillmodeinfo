// +build wasm,web

package modal

import (
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"syscall/js"
)

var modalStatus bool
var modalChan chan struct{}
var modalHidden = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	close(modalChan)
	js.Global().Call("$", "#theModal").Call("modal", "dispose")
	return nil
})

func Verify(title string, text string, button string) bool {
	document.GetElementById("modal-title").Set("innerText", title)
	document.GetElementById("modal-text").Set("innerText", text)
	document.GetElementById("modal-yes").Set("innerText", button)

	modalStatus = false
	modalChan=make(chan struct{})
	js.Global().Call("$", "#theModal").Call("modal")
	js.Global().Call("$", "#theModal").Call("on", "hidden.bs.modal", modalHidden)
	<-modalChan
	return modalStatus
}

func init() {
	js.Global().Set("cmiModalClick", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		modalStatus=true
		js.Global().Call("$", "#theModal").Call("modal", "hide")
		return nil
	}))
}