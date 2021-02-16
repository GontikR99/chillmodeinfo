// +build wasm

package jsbinding

import "syscall/js"

func OneShot(callback func(this js.Value, args []js.Value)interface{}) js.Func{
	funcRes:=new(js.Func)
	*funcRes=js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		funcRes.Release()
		return callback(this, args)
	})
	return *funcRes
}