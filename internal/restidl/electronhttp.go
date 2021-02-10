// +build wasm,electron

package restidl

import (
	"bytes"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"github.com/GontikR99/chillmodeinfo/pkg/jsbinding"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"syscall/js"
)

var https = nodejs.Require("https")

func httpCall(method string, path string, reqText []byte) (resBody []byte, statCode int, err error) {
	options := map[string]interface{}{
		"hostname": sitedef.DNSName,
		"port":     sitedef.Port,
		"path":     path,
		"method":   method,
		"headers": map[string]interface{}{
			"Content-Type":   "application/json",
			"Accept":         "application/json",
			"Content-Length": len(reqText),
		},
	}

	doneChan := make(chan struct{})
	buffer := new(bytes.Buffer)
	errHolder := new(error)
	statusCodeHolder := new(int)

	responseFunc := new(js.Func)
	*responseFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		responseFunc.Release()
		res := args[0]
		*statusCodeHolder = res.Get("statusCode").Int()

		onDataFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			chunk := args[0]
			buffer.Write(jsbinding.ReadArrayBuffer(chunk))
			return nil
		})
		onEndFunc := new(js.Func)
		onErrorFunc := new(js.Func)
		*onEndFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onDataFunc.Release()
			onEndFunc.Release()
			onErrorFunc.Release()
			doneChan <- struct{}{}
			return nil
		})
		*onErrorFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			onDataFunc.Release()
			onEndFunc.Release()
			onErrorFunc.Release()
			*errHolder = errors.New(args[0].String())
			return nil
		})
		res.Call("on", "data", onDataFunc)
		res.Call("on", "error", *onErrorFunc)
		res.Call("on", "end", *onEndFunc)

		return nil
	})

	req := https.Call("request", options, responseFunc)
	req.Call("write", jsbinding.BufferOf(reqText))
	req.Call("end")

	<-doneChan
	if *errHolder != nil {
		return nil, 0, *errHolder
	} else {
		return buffer.Bytes(), *statusCodeHolder, nil
	}
}
