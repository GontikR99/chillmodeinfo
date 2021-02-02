// +build wasm

package browserwindow

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcmain"
	"io"
	"syscall/js"
)

var browserWindow = electron.Get().Get("BrowserWindow")

type BrowserWindow interface {
	io.Closer
	ipc.Endpoint

	RemoveMenu()
	Show()
	LoadFile(path string)
	On(eventName string, action func())
	Once(eventName string, action func())
	SetAlwaysOnTop(bool)
}

type electronBrowserWindow struct {
	browserWindow js.Value
	webContents js.Value

	callbacks     map[int]js.Func
	nextCallback  int
}

type WebPreferences struct {
	Preload interface{} `json:"preload"`
	ContextIsolation interface{} `json:"contextIsolation"`
	NodeIntegration interface{} `json:"nodeIntegration"`
}

type Conf struct {
	Width       interface{}  `json:"width"`
	Height      interface{}  `json:"height"`
	Show        interface{} `json:"show"`
	Transparent interface{} `json:"transparent"`
	Frame       interface{} `json:"frame"`
	WebPreferences *WebPreferences `json:"webPreferences"`
}

func trimPrefs(prefMap *map[string]interface{}) {
	var trimKeys []string
	for k,v := range *prefMap {
		if v==nil {
			trimKeys = append(trimKeys, k)
		}
		if submap, ok :=v.(map[string]interface{}); ok {
			trimPrefs(&submap)
			(*prefMap)[k]=submap
		}
	}

	for _, k := range trimKeys {
		delete(*prefMap, k)
	}
}

func New(conf Conf) BrowserWindow {
	data, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	parsed:=make(map[string]interface{})
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		panic(err)
	}
	trimPrefs(&parsed)
	browserWindowInstance :=browserWindow.New(parsed)

	return &electronBrowserWindow{
		browserWindow: browserWindowInstance,
		webContents: browserWindowInstance.Get("webContents"),
		callbacks:     make(map[int]js.Func),
		nextCallback:  0,
	}
}

func (bw *electronBrowserWindow) registerCallback(callback func()) (int, js.Func) {
	wrapped := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback()
		return nil
	})
	bw.nextCallback++
	bw.callbacks[bw.nextCallback] = wrapped
	return bw.nextCallback, wrapped
}

func (bw *electronBrowserWindow) singleshotCallback(callback func()) js.Func {
	var cbHolder int
	cbLoc := &cbHolder
	cbId, wrapped := bw.registerCallback(func() {
		callback()
		if fnc, ok := bw.callbacks[*cbLoc]; ok {
			fnc.Release()
			delete(bw.callbacks, *cbLoc)
		}
	})
	*cbLoc = cbId
	return wrapped
}

func (bw *electronBrowserWindow) Close() error {
	for _, w := range bw.callbacks {
		w.Release()
	}
	bw.callbacks = make(map[int]js.Func)
	return nil
}

func (bw *electronBrowserWindow) RemoveMenu() {
	bw.browserWindow.Call("removeMenu")
}

func (bw *electronBrowserWindow) Show() {
	bw.browserWindow.Call("show")
}

func (bw *electronBrowserWindow) LoadFile(path string) {
	bw.browserWindow.Call("loadFile", path)
}

func (bw *electronBrowserWindow) Once(eventName string, action func()) {
	bw.browserWindow.Call("once", eventName, bw.singleshotCallback(action))
}

func (bw *electronBrowserWindow) On(eventName string, action func()) {
	_, wrapped := bw.registerCallback(action)
	bw.browserWindow.Call("on", eventName, wrapped)
}

func (bw *electronBrowserWindow) SetAlwaysOnTop(b bool) {
	bw.browserWindow.Call("setAlwaysOnTop", b)
}

func (bw *electronBrowserWindow) Send(channelName string, content []byte) {
	bw.webContents.Call("send", ipc.Prefix+channelName, string(content))
}

func (bw *electronBrowserWindow) Listen(channelName string) <-chan ipc.Message {
	return ipcmain.Listen(channelName)
}
