// +build wasm,web

package ui

import "github.com/vugu/vugu"

type ChangeEvent interface {
	Value() string
	SetValue(string)
	Env() vugu.EventEnv
}
//vugugen:event Change
