// +build wasm,web

package settings

import (
	"github.com/GontikR99/chillmodeinfo/internal/console"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/vugu/vugu"
	"syscall/js"
)

type Settings struct {
	EqDir string
}

func (c *Settings) Rendered(_ vugu.RenderedCtx) {
	eqDirValue := js.Global().Get("document").Call("getElementById","settings-eqdir").Get("value").String()
	if eqDirValue != c.EqDir {
		js.Global().Get("document").Call("getElementById", "settings-eqdir").Set("value", c.EqDir)
	}
}

func (c *Settings) UpdateEqDir() {
	c.EqDir = js.Global().Get("document").Call("getElementById","settings-eqdir").Get("value").String()
	console.Log(c.EqDir)
}

func (c *Settings) BrowseEqDir(event vugu.DOMEvent) {
	event.PreventDefault()
	go func() {
		newDir, err := rpcidl.DirectoryDialog(ipcrenderer.Client, c.EqDir)
		if err!=nil {
			console.Log(err.Error())
			return
		}
		event.EventEnv().Lock()
		c.EqDir = newDir
		event.EventEnv().UnlockRender()
	}()
}