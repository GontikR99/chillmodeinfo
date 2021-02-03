// +build wasm,electron

package exerpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
)

var dialog=electron.Get().Get("dialog")

func init() {
	register(rpcidl.HandleDirectoryDialog(func(initial string)(string, error) {
		success, error :=nodejs.Promise(dialog.Call("showOpenDialog", map[string]interface{} {
			"title": "Select a directory",
			"defaultPath": initial,
			"properties": []interface{}{
				"openDirectory",
			},
		}))
		select {
			case successObj:=<-success:
				if successObj.Get("canceled").Bool() {
					return "", errors.New("Selection canceled")
				} else {
					filePaths:=successObj.Get("filePaths")
					if filePaths.Length()!=1 {
						return "", errors.New("Single selection expected")
					}
					return filePaths.Index(0).String(), nil
				}
			case errMsg:=<-error:
				return "", errors.New(errMsg.String())
		}
	}))
}