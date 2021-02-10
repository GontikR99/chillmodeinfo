// +build wasm,web

package signins

import (
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/document"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"sync"
	"syscall/js"
)

var googleReady = &sync.WaitGroup{}
var gapi js.Value
var auth2 js.Value

func init() {
	if !electron.IsPresent() {
		googleReady.Add(1)

		auth2ReadyFunc := new(js.Func)
		*auth2ReadyFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			auth2ReadyFunc.Release()
			auth2=gapi.Get("auth2").Call("init", map[string]interface{}{
				"client_id":    sitedef.GoogleSigninClientId,
				"cookiepolicy": "single_host_origin",
			})

			auth2.Get("currentUser").Call("listen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				SetToken(TokenGoogle + args[0].Call("getAuthResponse").Get("id_token").String())
				return nil
			}))

			auth2.Get("isSignedIn").Call("listen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				if !args[0].Bool() {
					ClearToken()
				}
				return nil
			}))

			if auth2.Get("isSignedIn").Call("get").Bool() {
				auth2.Call("signIn")
			}
			googleReady.Done()
			return nil
		})

		loginReadyFunc := new(js.Func)
		*loginReadyFunc = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
			loginReadyFunc.Release()
			gapi=js.Global().Get("gapi")
			gapi.Call( "load", "auth2", auth2ReadyFunc)

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

func getAuth2() js.Value {
	googleReady.Wait()
	return auth2
}

func SignIn() {
	if electron.IsPresent() {
		go rpcidl.GetSignIn(ipcrenderer.Client).SignIn()
	} else {
		go func() {
			getAuth2().Call("signIn")
		}()
	}
}

func SignOut() {
	if electron.IsPresent() {
		go rpcidl.GetSignIn(ipcrenderer.Client).SignOut()
	} else {
		go func() {
			getAuth2().Call("signOut")
		}()
	}
}
