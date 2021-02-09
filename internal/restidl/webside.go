// +build wasm,web

package restidl

import (
	"encoding/json"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/pkg/jsbinding"
	"github.com/GontikR99/chillmodeinfo/pkg/signin"
	"net/http"
	"syscall/js"
)

func call(method string, path string, request interface{}, response interface{}) error {
	requestWrapper := &packagedRequest{ReqMsg: request}
	if electron.IsPresent() {
		clientId, present, err := rpcidl.LookupSetting(ipcrenderer.Client, settings.ClientId)
		if err != nil {
			return err
		}
		if !present {
			return errors.New("No clientId generated")
		}
		requestWrapper.ClientId = clientId
	} else {
		curSignIn := signin.CurrentSignIn()
		if curSignIn != nil {
			requestWrapper.IdToken = curSignIn.IdToken()
		}
	}

	reqText, err := json.Marshal(requestWrapper)
	if err != nil {
		return err
	}

	doneChan := make(chan struct{})
	xhr := js.Global().Get("XMLHttpRequest").New()
	readyStateChange := new(js.Func)
	*readyStateChange = js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		if xhr.Get("readyState").Int() == 4 {
			// XMLHttpRequest.DONE
			readyStateChange.Release()
			close(doneChan)
		}
		return nil
	})
	xhr.Call("addEventListener", "readystatechange", *readyStateChange)
	if electron.IsPresent() {
		xhr.Call("open", method, sitedef.SiteURL+path, true)
	} else {
		xhr.Call("open", method, path, true)
	}
	xhr.Set("responseType", "arraybuffer")
	xhr.Call("setRequestHeader", "Accept", "application/json")
	xhr.Call("send", jsbinding.MakeArrayBuffer(reqText))
	<-doneChan
	resBody := jsbinding.ReadArrayBuffer(xhr.Get("response"))
	if resBody == nil {
		return errors.New("Failed to recover payload")
	}
	resCode := xhr.Get("status").Int()
	if resCode == http.StatusForbidden {
		signin.SignOut()
		err = errors.New(string(resBody))
		toast.Error("identity", err)
		return err
	}
	if resCode >= 400 {
		return errors.New(string(resBody))
	}
	return json.Unmarshal(resBody, response)
}

// Do nothing call so that IDL files can build client side
func serve(mux *http.ServeMux, path string, handler func(method string, request *Request) (interface{}, error)) {
}
