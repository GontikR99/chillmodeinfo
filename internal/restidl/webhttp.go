// +build wasm,web

package restidl

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/jsbinding"
	"syscall/js"
)

func httpCall(method string, path string, reqText []byte) (resBody []byte, statCode int, err error) {
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
	resBody = jsbinding.ReadArrayBuffer(xhr.Get("response"))
	if resBody == nil {
		err = errors.New("Failed to recover payload")
		return
	}
	statCode = xhr.Get("status").Int()
	return
}
