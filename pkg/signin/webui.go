// +build wasm,web

package signin

import (
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"syscall/js"
)

func RenderSignin(elementId string) {
	if !electron.IsPresent() {
		onSuccess := new(js.Func)
		*onSuccess = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onSuccess.Release()
			user := args[0]
			profile := user.Call("getBasicProfile")
			currentSignIn = &UserInfo{
				ID:      profile.Call("getId").String(),
				Name:    profile.Call("getName").String(),
				Email:   profile.Call("getEmail").String(),
				jsValue: user,
			}
			invokeHandlers()
			return nil
		})
		onFailure := new(js.Func)
		*onFailure = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onFailure.Release()
			console.Log(args[0])
			return nil
		})
		go func() {
			waitForGoogleReady()
			js.Global().Get("gapi").Get("signin2").Call("render", elementId, map[string]interface{}{
				"scope":     "profile email",
				"width":     "160",
				"height":    "30",
				"longtitle": "true",
				"theme":     "dark",
				"onsuccess": *onSuccess,
				"onfailure": *onFailure,
			})
		}()
	}
}

func SignOut() {
	if !electron.IsPresent() {
		go func() {
			waitForGoogleReady()
			auth2 := js.Global().Get("gapi").Get("auth2").Call("getAuthInstance")
			successChan, errChan := nodejs.FromPromise(auth2.Call("signOut"))
			select {
			case <-successChan:
				currentSignIn = nil
				invokeHandlers()
			case errVal := <-errChan:
				console.LogRaw(errVal)
			}
		}()
	}
}
