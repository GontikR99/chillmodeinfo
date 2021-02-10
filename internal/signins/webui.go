// +build wasm,web

package signins

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/dom"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"syscall/js"
)

func Attach(signinButton dom.Element) {
	if electron.IsPresent() {
		toast.Error("coding", errors.New("reached signins.Attach within Electron"))
	} else {
		go func() {
			waitForGoogleReady()
			onSuccess := new(js.Func)
			*onSuccess = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				SetToken(TokenGoogle + args[0].Call("getAuthResponse").Get("id_token").String())
				return nil
			})
			onFailure := new(js.Func)
			*onFailure = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
				return nil
			})
			auth2 := js.Global().Get("gapi").Get("auth2")
			auth2.Call("init", map[string]interface{}{
				"client_id":    sitedef.GoogleSigninClientId,
				"cookiepolicy": "single_host_origin",
			})

			auth2Instance := auth2.Call("getAuthInstance")
			auth2Instance.Call("attachClickHandler", signinButton.JSValue(), map[string]interface{}{}, *onSuccess, *onFailure)
		}()
	}
}

func SignOut() {
	if electron.IsPresent() {
		go rpcidl.GetSignIn(ipcrenderer.Client).SignOut()
	} else {
		go func() {
			waitForGoogleReady()
			auth2 := js.Global().Get("gapi").Get("auth2").Call("getAuthInstance")
			successChan, errChan := nodejs.FromPromise(auth2.Call("signOut"))
			select {
			case <-successChan:
				ClearToken()
			case errVal := <-errChan:
				console.LogRaw(errVal)
			}
		}()
	}
}
