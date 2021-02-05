// +build wasm

package main

import (
"github.com/GontikR99/chillmodeinfo/pkg/console"
"github.com/vugu/vugu"
"github.com/vugu/vugu/domrender"
"github.com/vugu/vugu/js"
)

func main() {
	renderer, err := domrender.New("#page_root")
	if err != nil {
		panic(err)
	}
	defer renderer.Release()

	buildEnv, err := vugu.NewBuildEnv(renderer.EventEnv())
	if err != nil {
		panic(err)
	}

	root := &Root{}

	for ok := true; ok; ok = renderer.EventWait() {
		buildResults := buildEnv.RunBuild(root)

		err = renderer.Render(buildResults)
		js.Global().Get("feather").Call("replace")

		if err != nil {
			console.Log(err)
		}
	}
}

