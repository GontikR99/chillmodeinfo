// +build wasm,web

package util

import "syscall/js"

var toLocaleString=js.Global().Get("Number").Get("prototype").Get("toLocaleString")

func FormatFloat(value float64) string {
	return toLocaleString.Call("call", value).String()
}