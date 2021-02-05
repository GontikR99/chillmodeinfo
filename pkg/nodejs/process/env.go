// +build wasm,electron

package process

import "syscall/js"

var process=js.Global().Get("process")
var env= process.Get("env")

func LookupEnv(key string) (string, bool) {
	v := env.Get(key)
	if v.Type() == js.TypeString {
		return v.String(), true
	} else {
		return "", false
	}
}

func Getenv(key string) string {
	v, _ := LookupEnv(key)
	return v
}