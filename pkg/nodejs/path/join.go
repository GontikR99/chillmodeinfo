// +build wasm,electron

package path

import "github.com/GontikR99/chillmodeinfo/pkg/nodejs"

var path= nodejs.Require("path")

func Join(args ...string) string {
	var argArray []interface{}
	for _, arg := range args {
		argArray = append(argArray, arg)
	}
	return path.Call("join", argArray...).String()
}