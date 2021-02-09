// +build wasm,web

package signin

import (
	"github.com/GontikR99/chillmodeinfo/pkg/document"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"sync"
	"syscall/js"
)

var googleReady = &sync.WaitGroup{}

// Call at webapp startup.  Prepare for Google signin
func PrepareForSignin(clientId string) {
	if !electron.IsPresent() {
		googleReady.Add(1)

		metaTag := document.CreateElement("meta")
		metaTag.SetAttribute("name", "google-signin-client_id")
		metaTag.SetAttribute("content", clientId)
		head := document.GetElementsByTagName("head")[0]
		head.AppendChild(metaTag)

		loginReadyFunc := new(js.Func)
		*loginReadyFunc = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
			loginReadyFunc.Release()
			googleReady.Done()
			return nil
		})
		js.Global().Set("googleLibraryLoaded", *loginReadyFunc)

		scriptTag := document.CreateElement("script")
		scriptTag.SetAttribute("src", "https://apis.google.com/js/client:platform.js?onload=googleLibraryLoaded")
		scriptTag.SetAttribute("async", "")
		scriptTag.SetAttribute("defer", "")

		body := document.GetElementsByTagName("body")[0]
		body.AppendChild(scriptTag)
	}
}

func waitForGoogleReady() {
	googleReady.Wait()
}
