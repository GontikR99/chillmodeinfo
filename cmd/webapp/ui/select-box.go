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

//
//func (c *SelectBox) onChange(event js.Value, env vugu.EventEnv) {
//	newValue := event.Get("target").Get("value").String()
//	if c.Value != newValue {
//		c.Value = newValue
//		if c.Change != nil {
//			c.Change.ChangeHandle(&selectBoxChangeEvent{
//				value: c.Value,
//				env:   env,
//				sb:    c,
//			})
//		}
//	}
//}
//
//type selectBoxChangeEvent struct {
//	value string
//	env   vugu.EventEnv
//	sb    *SelectBox
//}
//
//func (s *selectBoxChangeEvent) Value() string      { return s.value }
//func (s *selectBoxChangeEvent) Env() vugu.EventEnv { return s.env }
//func (s *selectBoxChangeEvent) SetValue(s2 string) {
//	go func() {
//		s.env.Lock()
//		s.sb.Value = s2
//		s.env.UnlockRender()
//	}()
//}
//
//func (c *SelectBox) Init(vCtx vugu.InitCtx) {
//	c.ctx, c.doneFunc = context.WithCancel(context.Background())
//
//	if c.Id==0 {
//		selectBoxIdGen++
//		c.Id=selectBoxIdGen
//	}
//	c.idStr = fmt.Sprintf("select-box-%d", c.Id)
//	c.changeFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//		c.onChange(args[0], vCtx.EventEnv())
//		return nil
//	})
//	go func() {
//		for {
//			select {
//			case <-c.ctx.Done():
//				return
//			case <-time.After(10 * time.Millisecond):
//			}
//			elt := document.GetElementById(c.idStr)
//			if elt != nil {
//				elt.JSValue().Set("onchange", c.changeFunc)
//				if v:=elt.JSValue().Get("value").String(); v!=c.Value && c.Options!=nil {
//					selectedIdx:=-1
//					for idx, value := range c.Options {
//						if c.Value==value {
//							selectedIdx=idx
//							break
//						}
//					}
//					if selectedIdx!=-1 {
//						for idx, value := range c.Options {
//							elt := document.GetElementById(fmt.Sprintf("%s-option-%d", c.idStr, idx))
//							if elt != nil {
//								if c.Value == value {
//									elt.SetAttribute("selected", "selected")
//								} else {
//									elt.RemoveAttribute("selected")
//								}
//							}
//						}
//					}
//				}
//			}
//		}
//	}()
//}
//
//func (c *SelectBox) Destroy(vCtx vugu.DestroyCtx) {
//	c.doneFunc()
//	c.changeFunc.Release()
//}
