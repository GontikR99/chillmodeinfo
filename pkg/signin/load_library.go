// +build wasm,web

package signin

import (
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"sync"
	"syscall/js"
)

var googleReady = &sync.WaitGroup{}

// Call at webapp startup.  Prepare for Google signin
func PrepareForSignin(clientId string) {
	if !electron.IsPresent() {
		googleReady.Add(1)

		metaTag := js.Global().Get("document").Call("createElement", "meta")
		metaTag.Call("setAttribute", "name", "google-signin-client_id")
		metaTag.Call("setAttribute", "content", clientId)
		head := js.Global().Get("document").Call("getElementsByTagName", "head").Index(0)
		head.Call("appendChild", metaTag)

		loginReadyFunc := new(js.Func)
		*loginReadyFunc = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
			loginReadyFunc.Release()
			googleReady.Done()
			return nil
		})
		js.Global().Set("googleLibraryLoaded", *loginReadyFunc)

		scriptTag := js.Global().Get("document").Call("createElement", "script")
		scriptTag.Call("setAttribute", "src", "https://apis.google.com/js/client:platform.js?onload=googleLibraryLoaded")
		scriptTag.Call("setAttribute", "async", "")
		scriptTag.Call("setAttribute", "defer", "")

		body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
		body.Call("appendChild", scriptTag)
	}
}

func waitForGoogleReady() {
	googleReady.Wait()
}
