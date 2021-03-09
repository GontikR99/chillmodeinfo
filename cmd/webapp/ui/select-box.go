// +build wasm,web

package ui

import (
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/vuguutil"
	"github.com/vugu/vugu"
	"syscall/js"
)

var selectBoxIdGen int

type SelectBox struct {
	vuguutil.BackgroundComponent
	AttrMap vugu.AttrMap
	Value   string
	Change  ChangeHandler
	Options []string

	env   vugu.EventEnv
	idStr string
}

func (c *SelectBox) optNodeId(idx int) string {
	return fmt.Sprintf("%s-option-%d", c.idStr, idx)
}

func (c *SelectBox) Init(vCtx vugu.InitCtx) {
	selectBoxIdGen++
	c.idStr = fmt.Sprintf("selectbox-%d", selectBoxIdGen)

	activeSelectBoxes[c.idStr] = c
	c.env = vCtx.EventEnv()

	c.InitBackground(vCtx, c)
	c.ListenForRender()
}

func (c *SelectBox) RunInBackground() {
	defer func() {
		delete(activeSelectBoxes, c.idStr)
	}()
	for {
		select {
		case <-c.Done():
			return
		case <-c.Rendered():
			for idx, optVal := range c.Options {
				optElt := vuguutil.GetElementByNodeId(c.optNodeId(idx))
				if optElt == nil {
					continue
				}
				eltVal := optElt.GetAttribute("value")
				if eltVal.IsNull() || eltVal.IsUndefined() || eltVal.String() != optVal || optVal != c.Value {
					optElt.RemoveAttribute("selected")
				} else {
					optElt.SetAttribute("selected", "selected")
				}
			}
		}
	}
}

type selectBoxChangeEvent struct {
	value string
	env   vugu.EventEnv
	box   *SelectBox
}

var activeSelectBoxes = make(map[string]*SelectBox)

const sbChangeFunc = "cmiUiSelectBoxChange"

func init() {
	js.Global().Set(sbChangeFunc, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		boxName := args[1].String()
		if box, ok := activeSelectBoxes[boxName]; ok {
			go func() {
				box.env.Lock()
				box.onChange(vuguutil.NewVuguEvent(event, box.env))
				box.env.UnlockRender()
			}()
		}
		return nil
	}))
}

func (c *SelectBox) onChangeHookText() string {
	return sbChangeFunc + "(event, \"" + c.idStr + "\")"
}

func (s *selectBoxChangeEvent) Value() string      { return s.value }
func (s *selectBoxChangeEvent) Env() vugu.EventEnv { return s.box.env }
func (s *selectBoxChangeEvent) SetValue(s2 string) { s.box.Value = s2 }

func (c *SelectBox) onChange(event vugu.DOMEvent) {
	oldValue := c.Value
	c.Value = event.PropString("target", "value")
	if c.Change != nil && c.Value != oldValue {
		c.Change.ChangeHandle(&selectBoxChangeEvent{
			value: c.Value,
			box:   c,
		})
	}
}
