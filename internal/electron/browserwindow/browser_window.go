package browserwindow

import (
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"github.com/GontikR99/chillmodeinfo/internal/nodejs"
	"io"
	"syscall/js"
)

var browserWindow = electron.Get().Get("BrowserWindow")
var nodePath = nodejs.Require("path")

type BrowserWindow interface {
	io.Closer
	RemoveMenu()
	Show()
	LoadFile(path string)
	On(eventName string, action func())
	Once(eventName string, action func())
	SetAlwaysOnTop(bool)
}

type electronBrowserWindow struct {
	browserWindow js.Value
	callbacks map[int]js.Func
	nextCallback int
}

type Conf struct {
	Width int `json:"width"`
	Height int `json:"height"`
	Show bool `json:"show"`
	Transparent bool `json:"transparent"`
	Frame bool `json:"frame"`
}

func NewConf() *Conf {
	return &Conf{
		Width:       1024,
		Height:      768,
		Show:        true,
		Transparent: false,
		Frame:       true,
	}
}

func (bwc *Conf) WithWidth(width int) *Conf {
	bwc.Width=width
	return bwc
}

func (bwc *Conf) WithHeight(height int) *Conf {
	bwc.Height=height
	return bwc
}

func (bwc *Conf) WithShow(show bool) *Conf {
	bwc.Show=show
	return bwc
}

func (bwc *Conf) WithTransparent(transparent bool) *Conf {
	bwc.Transparent=transparent
	return bwc
}

func (bwc *Conf) WithFrame(frame bool) *Conf {
	bwc.Frame=frame
	return bwc
}

func New(conf *Conf) BrowserWindow {
	data, err := json.Marshal(conf)
	if err!=nil {
		panic(err)
	}
	jsv := js.Global().Get("JSON").Call("parse", string(data))
	return &electronBrowserWindow{
		browserWindow: browserWindow.New(jsv),
		callbacks:     make(map[int]js.Func),
		nextCallback:  0,
	}
}

func (bw *electronBrowserWindow) registerCallback(callback func()) (int, js.Func) {
	wrapped := js.FuncOf(func(this js.Value, args []js.Value)interface{} {
		callback()
		return nil
	})
	bw.nextCallback++
	bw.callbacks[bw.nextCallback]=wrapped
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
	*cbLoc=cbId
	return wrapped
}

func (bw *electronBrowserWindow) Close() error {
	for _, w := range bw.callbacks {
		w.Release()
	}
	bw.callbacks=make(map[int]js.Func)
	return nil
}

func (bw *electronBrowserWindow) RemoveMenu() {
	bw.browserWindow.Call("removeMenu")
}

func (bw *electronBrowserWindow) Show() {
	bw.browserWindow.Call("show")
}

func (bw *electronBrowserWindow) LoadFile(path string) {
	bw.browserWindow.Call("loadFile", nodePath.Call("join", electron.RootDirectory(), path))
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
