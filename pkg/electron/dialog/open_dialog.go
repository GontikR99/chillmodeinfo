// +build wasm,electron

package dialog

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/pkg/electron"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/binding"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/browserwindow"
	"github.com/GontikR99/chillmodeinfo/pkg/nodejs"
	"syscall/js"
)

var dialog= electron.JSValue().Get("dialog")

type OpenOptions struct {
	Title interface{}       `json:"title"`
	DefaultPath interface{} `json:"defaultPath"`
	ButtonLabel interface{} `json:"buttonLabel"`
	Filters *[]FileFilter   `json:"filters"`
	Properties *[]string    `json:"properties"`
}

type FileFilter struct {
	Name string `json:"name"`
	Extensions []string `json:"extensions"`
}

const (
	OpenFile = "openFile"
	OpenDirectory = "openDirectory"
	MultiSelections = "multiSelections"
	ShowHiddenFiles = "showHiddenFiles"
	PromptToCreate = "promptToCreate"
	DontAddToRecent = "dontAddToRecent"
)

func ShowOpenDialogModal(bw browserwindow.BrowserWindow, options *OpenOptions) ([]string, error) {
	var promiseval js.Value
	jsonOptions := binding.JsonifyOptions(options)
	if bw==nil {
		promiseval = dialog.Call("showOpenDialog", jsonOptions)
	} else {
		promiseval = dialog.Call("showOpenDialog", bw.JSValue(), jsonOptions)
	}
	successChan, errorChan := nodejs.FromPromise(promiseval)
	select {
	case err:=<-errorChan:
		return nil, errors.New(err[0].String())
	case succobj:=<-successChan:
		if succobj[0].Get("canceled").Bool() {
			return nil, errors.New("open canceled")
		}
		var filePaths []string
		filePathsJS := succobj[0].Get("filePaths")
		for i:=0;i<filePathsJS.Length();i++ {
			filePaths=append(filePaths, filePathsJS.Index(i).String())
		}
		return filePaths, nil
	}
}

func ShowOpenDialog(options *OpenOptions) ([]string, error) {
	return ShowOpenDialogModal(nil, options)
}