// +build wasm,web

package settings

import (
	"github.com/GontikR99/chillmodeinfo/internal/console"
	"github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcrenderer"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/vugu/vugu"
)

type Settings struct {
	EqDir *EqDirValue
}

func (c *Settings) Init(ctx vugu.InitCtx) {
	c.EqDir=&EqDirValue{}
	go func() {
		value, present, err := rpcidl.LookupSetting(ipcrenderer.Client, settings.EverQuestDirectory)
		if err==nil && present {
			ctx.EventEnv().Lock()
			c.EqDir.Path=value
			ctx.EventEnv().UnlockRender()
		}
	}()
}

func (c *Settings) BrowseEqDir(event vugu.DOMEvent) {
	event.PreventDefault()
	go func() {
		newDir, err := rpcidl.DirectoryDialog(ipcrenderer.Client, c.EqDir.Path)
		if err != nil {
			console.Log(err.Error())
			return
		}
		event.EventEnv().Lock()
		c.EqDir.SetStringValue(newDir)
		event.EventEnv().UnlockRender()
	}()
}

type EqDirValue struct {
	Path string
}

func (e *EqDirValue) StringValue() string {
	return e.Path
}

func (e *EqDirValue) SetStringValue(s string) {
	e.Path=s
	go rpcidl.SetSetting(ipcrenderer.Client, settings.EverQuestDirectory, s)
}
