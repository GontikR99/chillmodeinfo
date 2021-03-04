// +build wasm,web

package settings

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
	"strings"
)

type Settings struct {
	EqDir    *ConfiguredValue
	BidStart *ConfiguredValue
	BidEnd   *ConfiguredValue
	BidClose *ConfiguredValue

	BadMath *ConfiguredValue
	allowNamed bool
}

func (c *Settings) Init(ctx vugu.InitCtx) {
	c.EqDir = &ConfiguredValue{
		Key:      settings.EverQuestDirectory,
		Callback: func(s string) { rpcidl.RestartScan(ipcrenderer.Client) },
	}
	c.BidStart = &ConfiguredValue{Key: settings.BidStartPattern}
	c.BidEnd = &ConfiguredValue{Key: settings.BidEndPattern}
	c.BidClose = &ConfiguredValue{Key: settings.BidClosePattern}
	c.BadMath = &ConfiguredValue{Key: settings.BadMathThreshold}

	c.EqDir.Init(ctx)
	c.BidStart.Init(ctx)
	c.BidEnd.Init(ctx)
	c.BidClose.Init(ctx)
	c.BadMath.Init(ctx)

	go func() {
		namedStr, present, err := rpcidl.LookupSetting(ipcrenderer.Client, settings.AllowNamedBids)
		if present && err==nil {
			ctx.EventEnv().Lock()
			if strings.EqualFold("true", namedStr) {
				c.allowNamed=true
			} else {
				c.allowNamed=false
			}
			ctx.EventEnv().UnlockRender()
		}
	}()
}

func (c *Settings) BrowseEqDir(event vugu.DOMEvent) {
	event.PreventDefault()
	go func() {
		newDir, err := rpcidl.DirectoryDialog(ipcrenderer.Client, c.EqDir.Value)
		if err != nil {
			console.Log(err.Error())
			return
		}
		event.EventEnv().Lock()
		c.EqDir.SetStringValue(newDir)
		event.EventEnv().UnlockRender()
	}()
}

type ConfiguredValue struct {
	Key      string
	Value    string
	Callback func(value string)
}

func (cv *ConfiguredValue) Init(ctx vugu.InitCtx) {
	go func() {
		value, present, err := rpcidl.LookupSetting(ipcrenderer.Client, cv.Key)
		if err == nil && present {
			ctx.EventEnv().Lock()
			cv.Value = value
			ctx.EventEnv().UnlockRender()
		}
	}()
}

func (cv *ConfiguredValue) StringValue() string {
	return cv.Value
}

func (cv *ConfiguredValue) SetStringValue(s string) {
	cv.Value = s
	go func() {
		rpcidl.SetSetting(ipcrenderer.Client, cv.Key, s)
		if cv.Callback != nil {
			cv.Callback(s)
		}
	}()
}

func (c *Settings) PositionOverlay(event vugu.DOMEvent, name string) {
	event.PreventDefault()
	go rpcidl.PositionOverlay(ipcrenderer.Client, name)
}

func (c *Settings) ResetOverlay(event vugu.DOMEvent, name string) {
	event.PreventDefault()
	go rpcidl.ResetOverlay(ipcrenderer.Client, name)
}

func (c *Settings) CloseOverlay(event vugu.DOMEvent, name string) {
	event.PreventDefault()
	go rpcidl.CloseOverlay(ipcrenderer.Client, name)
}

func (c *Settings) changeNamedBids(event vugu.DOMEvent) {
	c.allowNamed=event.PropBool("target", "checked")
	go func() {
		if c.allowNamed {
			rpcidl.SetSetting(ipcrenderer.Client, settings.AllowNamedBids, "true")
		} else {
			rpcidl.SetSetting(ipcrenderer.Client, settings.AllowNamedBids, "false")
		}
	}()
}