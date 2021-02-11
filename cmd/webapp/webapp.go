// +build wasm, web

package main

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/signins"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/domrender"
	"github.com/vugu/vugu/js"
	"log"
)

func main() {
	toast.ListenForToasts()
	signins.OnStateChange(func() {
		if signins.SignedIn() {
			err := restidl.VerifyLogin()
			if err != nil {
				signins.SignOut()
				toast.Error("accounts", errors.New("Unable to verify login: "+err.Error()))
			}
		}
	})
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
			log.Println(err)
		}
	}
}
