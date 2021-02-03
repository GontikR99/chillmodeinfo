// +build wasm,electron

package browserwindow

import (
	"encoding/hex"
	"encoding/json"
	"github.com/GontikR99/chillmodeinfo/internal/electron"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcmain"
	"github.com/GontikR99/chillmodeinfo/internal/msgcomm"
	"io"
	"net/rpc"
	"strconv"
	"syscall/js"
)

var browserWindow = electron.Get().Get("BrowserWindow")

type BrowserWindow interface {
	io.Closer
	msgcomm.Endpoint

	RemoveMenu()
	Show()
	LoadFile(path string)

	Id() int

	On(eventName string, action func())
	Once(eventName string, action func())

	SetAlwaysOnTop(bool)

	// Additions
	OnClosed(action func())
	ServeRPC(server *rpc.Server)
	JSValue() js.Value
}

type electronBrowserWindow struct {
	browserWindow js.Value
	webContents js.Value

	closedCallbacks []func()
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
	browserWindowInternal :=browserWindow.New(parsed)

	browserWindowInstance := &electronBrowserWindow{
		browserWindow: browserWindowInternal,
		webContents:   browserWindowInternal.Get("webContents"),
		callbacks:     make(map[int]js.Func),
		nextCallback:  0,
	}

	var handleClosedFunc js.Func
	handleClosedFuncAddr :=&handleClosedFunc
	*handleClosedFuncAddr = js.FuncOf(func(_ js.Value, _ []js.Value)interface{} {
		browserWindowInstance.handleClosed()
		(*handleClosedFuncAddr).Release()
		return nil
	})
	browserWindowInternal.Call("on", "closed", handleClosedFunc)

	return browserWindowInstance
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

// On window closed, call all of our queued close callbacks, then release any outstanding callbacks
func (bw *electronBrowserWindow) handleClosed() {
	for i:=len(bw.closedCallbacks)-1; i>=0;i-- {
		bw.closedCallbacks[i]()
	}
	for _, w := range bw.callbacks {
		w.Release()
	}
	bw.callbacks = make(map[int]js.Func)
}

// Do something when the window is closed
func (bw *electronBrowserWindow) OnClosed(callback func()) {
	bw.closedCallbacks = append(bw.closedCallbacks, callback)
}

// Close the window, as if the user had clicked the close button.  May be intercepted/interrupted by code
func (bw *electronBrowserWindow) Close() error {
	bw.browserWindow.Call("close")
	return nil
}

// Turn off the menu bar
func (bw *electronBrowserWindow) RemoveMenu() {
	bw.browserWindow.Call("removeMenu")
}

// Show the window
func (bw *electronBrowserWindow) Show() {
	bw.browserWindow.Call("show")
}

// Load content into the window
func (bw *electronBrowserWindow) LoadFile(path string) {
	bw.browserWindow.Call("loadFile", path)
}

// Register a callback to be called once
func (bw *electronBrowserWindow) Once(eventName string, action func()) {
	if eventName=="closed" {
		bw.closedCallbacks=append(bw.closedCallbacks, action)
	} else {
		bw.browserWindow.Call("once", eventName, bw.singleshotCallback(action))
	}
}

// Register a callback to be called repeatedly
func (bw *electronBrowserWindow) On(eventName string, action func()) {
	if eventName=="closed" {
		bw.closedCallbacks=append(bw.closedCallbacks, action)
	} else {
		_, wrapped := bw.registerCallback(action)
		bw.browserWindow.Call("on", eventName, wrapped)
	}
}

func (bw *electronBrowserWindow) SetAlwaysOnTop(b bool) {
	bw.browserWindow.Call("setAlwaysOnTop", b)
}

func (bw *electronBrowserWindow) Send(channelName string, content []byte) {
	bw.webContents.Call("send", msgcomm.Prefix+channelName, hex.EncodeToString(content))
}

func (bw *electronBrowserWindow) Id() int {
	return bw.browserWindow.Get("id").Int()
}

func (bw *electronBrowserWindow) Listen(channelName string) (<-chan msgcomm.Message, func()) {
	outChan := make(chan msgcomm.Message)
	inChan, inDone := ipcmain.Listen(channelName)
	go func() {
		for {
			select {
			case inMsg := <- inChan:
				if inMsg==nil {
					return
				}
				if inMsg.Sender() == strconv.Itoa(bw.Id()) {
					outChan <- inMsg
				}
			}
		}
	}()
	return outChan, inDone
}

func (bw *electronBrowserWindow) ServeRPC(server *rpc.Server) {
	endpointStream := msgcomm.EndpointAsStream("rpcMain", bw)
	bw.OnClosed(func() {
		endpointStream.Close()
	})
	go func() {
		server.ServeConn(endpointStream)
	}()
}

func (bw *electronBrowserWindow) JSValue() js.Value {
	return bw.browserWindow
}