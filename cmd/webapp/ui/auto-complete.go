// +build wasm,web

package ui

import (
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/vuguutil"
	"github.com/vugu/vugu"
	"time"
)

var autocompleteIdGen = 0

type AutoComplete struct {
	vuguutil.BackgroundComponent
	AttrMap    vugu.AttrMap
	Value      string
	Suggestion SuggestionHandler
	Change     ChangeHandler
	EditStyle  string

	editIdStr    string
	suggestIdStr string

	focused   chan vugu.DOMEvent
	blurred   chan vugu.DOMEvent
	keyDowned chan vugu.DOMEvent

	sugStyle string

	lastSetValue   string
	proposedValues []string
	suggestionIdx  int
	unselectText   string

	env vugu.EventEnv
}

func (c *AutoComplete) Init(vCtx vugu.InitCtx) {
	autocompleteIdGen++
	c.editIdStr = fmt.Sprintf("autocomplete-edit-%d", autocompleteIdGen)
	c.suggestIdStr = fmt.Sprintf("autocomplete-suggest-%d", autocompleteIdGen)
	c.env = vCtx.EventEnv()

	c.focused = make(chan vugu.DOMEvent)
	c.blurred = make(chan vugu.DOMEvent)
	c.keyDowned = make(chan vugu.DOMEvent)

	c.InitBackground(vCtx, c)
}

func (c *AutoComplete) editorText() string {
	editElt := vuguutil.GetElementByNodeId(c.editIdStr)
	if editElt == nil {
		return ""
	} else {
		return editElt.Get("value").String()
	}
}

func (c *AutoComplete) RunInBackground() {
	var event vugu.DOMEvent
	hasFocus := false
	for {
		editBoxElt := vuguutil.GetElementByNodeId(c.editIdStr)
		sugBoxElt := vuguutil.GetElementByNodeId(c.suggestIdStr)
		select {
		case <-c.Done():
			return
		case <-c.Rendered():
			if sugBoxElt != nil {
				sugBoxElt.SetAttribute("style", c.sugStyle)
			}
			if c.Value != c.lastSetValue {
				c.lastSetValue = c.Value
				if editBoxElt != nil {
					editBoxElt.Set("value", c.Value)
				}
				c.populateSuggestions(c.Value, c.env)
			}
		case event = <-c.focused:
			hasFocus = true
			event.JSEventTarget().Call("select")
			c.populateSuggestions(c.Value, c.env)

		case event = <-c.blurred:
			hasFocus = false
			c.populateSuggestions("", c.env)
			if editBoxElt != nil {
				eValue := editBoxElt.Get("value").String()
				if eValue != c.Value {
					c.Value = eValue
					c.lastSetValue = eValue
					if c.Change != nil {
						c.Change.ChangeHandle(&autocompleteChangeEvent{
							value: eValue,
							ac:    c,
						})
					}
				}
			}

		case event = <-c.keyDowned:
			keyCode := event.PropString("code")
			if keyCode == "ArrowUp" {
				event.PreventDefault()
				event.StopPropagation()
				go func() {
					c.env.Lock()
					defer c.env.UnlockRender()
					c.suggestionIdx--
					if c.suggestionIdx < -1 {
						c.suggestionIdx = -1
					}
					if editBoxElt != nil {
						if c.suggestionIdx == -1 {
							editBoxElt.Set("value", c.unselectText)
						} else {
							editBoxElt.Set("value", c.proposedValues[c.suggestionIdx])
						}
					}
				}()
			} else if keyCode == "ArrowDown" {
				event.PreventDefault()
				event.StopPropagation()
				go func() {
					c.env.Lock()
					defer c.env.UnlockRender()
					c.suggestionIdx++
					if c.suggestionIdx >= len(c.proposedValues) {
						c.suggestionIdx = len(c.proposedValues) - 1
					}
					if c.suggestionIdx>=0 && editBoxElt != nil {
						editBoxElt.Set("value", c.proposedValues[c.suggestionIdx])
					}
				}()
			} else {
				go func() {
					<-time.After(100 * time.Millisecond)
					if hasFocus {
						c.populateSuggestions(c.editorText(), c.env)
					}
				}()
			}
		case <-time.After(100 * time.Millisecond):
			if sugBoxElt == nil || editBoxElt == nil {
				break
			}
			if !hasFocus || c.proposedValues == nil {
				c.sugStyle = "position:absolute; display:none;"
			} else {
				c.sugStyle = fmt.Sprintf("position:absolute; "+
					"display:block; "+
					"border-style: solid; "+
					"border-width: 1px;"+
					"border-color: black; "+
					"margin: 2px; "+
					"color: black; "+
					"background-color:white; "+
					"top: %dpx; "+
					"left: -2px; "+
					"width: %dpx; "+
					"z-index: 100;",
					editBoxElt.JSValue().Get("offsetHeight").Int()-2,
					editBoxElt.JSValue().Get("offsetWidth").Int())
			}
			sugBoxElt.SetAttribute("style", c.sugStyle)
		}
	}
}

func (c *AutoComplete) onFocus(event vugu.DOMEvent) {
	c.focused <- event
}

func (c *AutoComplete) onBlur(event vugu.DOMEvent) {
	c.blurred <- event
}

func (c *AutoComplete) onKeyDown(event vugu.DOMEvent) {
	c.keyDowned <- event
}

func (c *AutoComplete) populateSuggestions(text string, env vugu.EventEnv) {
	go func() {
		if text == "" {
			env.Lock()
			c.proposedValues = nil
			c.suggestionIdx = -1
			env.UnlockRender()
			return
		}
		if c.Suggestion != nil {
			c.Suggestion.SuggestionHandle(&autocompleteSuggestionEvent{text, c})
		}
	}()
}

type autocompleteSuggestionEvent struct {
	value string
	ac    *AutoComplete
}

func (a *autocompleteSuggestionEvent) Value() string { return a.value }
func (a *autocompleteSuggestionEvent) Propose(strings []string) {
	if len(strings) > 20 {
		strings = strings[:20]
	}
	if a.value == a.ac.editorText() {
		a.ac.env.Lock()
		a.ac.unselectText = a.value
		a.ac.suggestionIdx = -1
		a.ac.proposedValues = strings
		a.ac.env.UnlockRender()
	}
}

type autocompleteChangeEvent struct {
	value string
	ac    *AutoComplete
}

func (a *autocompleteChangeEvent) Value() string      { return a.value }
func (a *autocompleteChangeEvent) SetValue(s string)  { a.ac.Value = s }
func (a *autocompleteChangeEvent) Env() vugu.EventEnv { return a.ac.env }

func (c *AutoComplete) suggestionStyle(idx int) string {
	if idx == c.suggestionIdx {
		return "width:100%; background-color: blue; color: white; cursor:pointer;"
	} else {
		return "width:100%; color: gray; cursor:pointer;"
	}
}

func (c *AutoComplete) suggestionClick(event vugu.DOMEvent, idx int) {
	event.PreventDefault()
	event.StopPropagation()
	c.suggestionIdx = idx
	input := vuguutil.GetElementByNodeId(c.editIdStr)
	if input != nil {
		input.Set("value", c.proposedValues[c.suggestionIdx])
	}
}

//
//func (c *AutoComplete) onFocus(event vugu.DOMEvent) {
//	if c.focusDone != nil {
//		c.focusDone()
//	}
//	var focusCtx context.Context
//	focusCtx, c.focusDone = context.WithCancel(c.ctx)
//	c.selectedSuggestion = -1
//	go func() {
//		defer func() {
//			suggestElt := document.GetElementById(c.suggestIdStr)
//			if suggestElt != nil {
//				suggestElt.SetAttribute("style", "position:absolute; display:none;")
//			}
//		}()
//		for {
//			select {
//			case <-focusCtx.Done():
//				return
//			case <-time.After(10 * time.Millisecond):
//			}
//			inputElt := document.GetElementById(c.editIdStr)
//			suggestElt := document.GetElementById(c.suggestIdStr)
//			if inputElt == nil || suggestElt == nil {
//				continue
//			}
//			inputElt.JSValue().Set("onkeydown", c.keyDownFunc)
//			if len(c.proposedValues) == 0 {
//				suggestElt.SetAttribute("style", "position:absolute; display:none;")
//			} else {
//				suggestElt.SetAttribute("style",
//					fmt.Sprintf("position:absolute; "+
//						"display:block; "+
//						"border-style: solid; "+
//						"border-width: 1px;"+
//						"border-color: black; "+
//						"margin: 2px; "+
//						"background-color:white; "+
//						"top: %dpx; "+
//						"left: -2px; "+
//						"width: %dpx; " +
//						"z-index: 100;",
//						inputElt.JSValue().Get("offsetHeight").Int()-2,
//						inputElt.JSValue().Get("offsetWidth").Int()))
//			}
//		}
//	}()
//}
//
//func (c *AutoComplete) onBlur(event vugu.DOMEvent) {
//	if c.focusDone != nil {
//		c.focusDone()
//		c.focusDone = nil
//	}
//	oldValue := c.Value
//	c.Value = c.displayValue
//	if c.Change != nil && oldValue != c.Value {
//		c.Change.ChangeHandle(&changeEvent{c.Value, c.env, c})
//	}
//}
//
//
//
//func (c *AutoComplete) onKeyDown(event js.Value, env vugu.EventEnv) {
//	code := event.Get("code").String()
//	switch code {
//	case "ArrowUp":
//		event.Call("preventDefault")
//		event.Call("stopPropagation")
//		c.selectedSuggestion--
//		if c.selectedSuggestion < -1 {
//			c.selectedSuggestion = -1
//		}
//		env.Lock()
//		if c.selectedSuggestion >= 0 {
//			event.Get("target").Set("value", c.proposedValues[c.selectedSuggestion])
//			event.Get("target").Call("select")
//			c.displayValue = c.proposedValues[c.selectedSuggestion]
//		}
//		env.UnlockRender()
//
//	case "ArrowDown":
//		event.Call("preventDefault")
//		event.Call("stopPropagation")
//		c.selectedSuggestion++
//		if c.selectedSuggestion >= len(c.proposedValues) {
//			c.selectedSuggestion = len(c.proposedValues) - 1
//			return
//		}
//		env.Lock()
//		if c.selectedSuggestion >= 0 {
//			event.Get("target").Set("value", c.proposedValues[c.selectedSuggestion])
//			event.Get("target").Call("select")
//			c.displayValue = c.proposedValues[c.selectedSuggestion]
//		}
//		env.UnlockRender()
//
//	default:
//		if c.selectedSuggestion != -1 {
//			env.Lock()
//			c.selectedSuggestion = -1
//			c.proposedValues = nil
//			env.UnlockRender()
//		}
//
//		go func() {
//			<-time.After(100 * time.Millisecond)
//			target := document.GetElementById(c.editIdStr)
//			if target != nil {
//				curValue := target.JSValue().Get("value").String()
//				c.displayValue = curValue
//				if len(curValue) == 0 {
//					env.Lock()
//					c.proposedValues = nil
//					env.UnlockRender()
//				} else if !c.queryPending && c.Suggestion != nil {
//					c.queryPending = true
//					c.Suggestion.SuggestionHandle(&suggestionEvent{
//						value: curValue,
//						env:   env,
//						ac:    c,
//					})
//				}
//			}
//		}()
//	}
//}
//
//type suggestionEvent struct {
//	value string
//	env   vugu.EventEnv
//	ac    *AutoComplete
//}
//
//func (s *suggestionEvent) Value() string { return s.value }
//
//func (s *suggestionEvent) Propose(strings []string) {
//	defer func() {
//		s.ac.queryPending = false
//	}()
//	target := document.GetElementById(s.ac.editIdStr)
//	if target != nil {
//		curValue := target.JSValue().Get("value").String()
//		if curValue != s.value {
//			return
//		}
//		if len(strings) > 20 {
//			strings = strings[:20]
//		}
//		s.env.Lock()
//		s.ac.proposedValues = strings
//		s.env.UnlockRender()
//	}
//}
//
//type changeEvent struct {
//	value string
//	env   vugu.EventEnv
//	ac    *AutoComplete
//}
//
//func (c *changeEvent) Value() string      { return c.value }
//func (c *changeEvent) Env() vugu.EventEnv { return c.env }
//func (c *changeEvent) SetValue(value string) {
//	c.ac.proposedValues = nil
//	c.ac.selectedSuggestion = -1
//	c.ac.Value = value
//	c.ac.displayValue = value
//	inputElt := document.GetElementById(c.ac.editIdStr)
//	if inputElt != nil {
//		inputElt.JSValue().Set("value", value)
//	}
//}
//
//func (c *AutoComplete) Init(vCtx vugu.InitCtx) {
//	c.ctx, c.doneFunc = context.WithCancel(context.Background())
//	c.displayValue = c.Value
//	c.env = vCtx.EventEnv()
//
//	if c.Id==0 {
//		autocompleteIdGen++
//		c.Id=autocompleteIdGen
//	}
//	c.editIdStr = fmt.Sprintf("autocomplete-input-%d", c.Id)
//	c.suggestIdStr = fmt.Sprintf("autocomplete-suggest-%d", c.Id)
//
//	c.keyDownFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//		c.onKeyDown(args[0], vCtx.EventEnv())
//		return nil
//	})
//}
//
//func (c *AutoComplete) Destroy(vCtx vugu.DestroyCtx) {
//	c.doneFunc()
//	c.keyDownFunc.Release()
//}
//
type SuggestionEvent interface {
	Value() string

	Propose([]string)
}

//vugugen:event Suggestion
