// +build wasm,web

package ui

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/vugu/vugu"
	"syscall/js"
	"time"
)

var autocompleteIdGen = 0

type AutoComplete struct {
	AttrMap    vugu.AttrMap
	Value      string
	Suggestion SuggestionHandler
	Change     ChangeHandler

	displayValue string

	editIdStr    string
	suggestIdStr string

	currentSuggestions []string
	selectedSuggestion int
	queryPending       bool

	ctx      context.Context
	doneFunc context.CancelFunc

	focusDone context.CancelFunc

	keyDownFunc js.Func
	env         vugu.EventEnv
}

func (c *AutoComplete) onFocus(event vugu.DOMEvent) {
	if c.focusDone != nil {
		c.focusDone()
	}
	var focusCtx context.Context
	focusCtx, c.focusDone = context.WithCancel(c.ctx)
	c.selectedSuggestion = -1
	go func() {
		defer func() {
			suggestElt := document.GetElementById(c.suggestIdStr)
			if suggestElt != nil {
				suggestElt.SetAttribute("style", "position:absolute; display:none;")
			}
		}()
		for {
			select {
			case <-focusCtx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
			inputElt := document.GetElementById(c.editIdStr)
			suggestElt := document.GetElementById(c.suggestIdStr)
			if inputElt == nil || suggestElt == nil {
				continue
			}
			inputElt.JSValue().Set("onkeydown", c.keyDownFunc)
			if len(c.currentSuggestions) == 0 {
				suggestElt.SetAttribute("style", "position:absolute; display:none;")
			} else {
				suggestElt.SetAttribute("style",
					fmt.Sprintf("position:absolute; "+
						"display:block; "+
						"border-style: solid; "+
						"border-width: 1px;"+
						"border-color: black; "+
						"margin: 2px; "+
						"background-color:white; "+
						"top: %dpx; "+
						"left: -2px; "+
						"width: %dpx; ",
						inputElt.JSValue().Get("offsetHeight").Int()-2,
						inputElt.JSValue().Get("offsetWidth").Int()))
			}
		}
	}()
}

func (c *AutoComplete) onBlur(event vugu.DOMEvent) {
	if c.focusDone != nil {
		c.focusDone()
		c.focusDone = nil
	}
	oldValue := c.Value
	c.Value = c.displayValue
	if c.Change != nil && oldValue != c.Value {
		c.Change.ChangeHandle(&changeEvent{c.Value, c.env, c})
	}
}

func (c *AutoComplete) suggestionStyle(idx int) string {
	if idx == c.selectedSuggestion {
		return "width:100%; background-color: blue; text-color: white; cursor:pointer;"
	} else {
		return "width:100%; text-color: gray; cursor:pointer;"
	}
}

func (c *AutoComplete) suggestionClick(event vugu.DOMEvent, idx int) {
	event.PreventDefault()
	event.StopPropagation()
	c.selectedSuggestion = idx
	input := document.GetElementById(c.editIdStr)
	if input != nil {
		input.JSValue().Set("value", c.currentSuggestions[c.selectedSuggestion])
		input.JSValue().Call("select")
	}
	c.displayValue = c.currentSuggestions[c.selectedSuggestion]
}

func (c *AutoComplete) onKeyDown(event js.Value, env vugu.EventEnv) {
	code := event.Get("code").String()
	switch code {
	case "ArrowUp":
		event.Call("preventDefault")
		event.Call("stopPropagation")
		c.selectedSuggestion--
		if c.selectedSuggestion < -1 {
			c.selectedSuggestion = -1
		}
		env.Lock()
		if c.selectedSuggestion >= 0 {
			event.Get("target").Set("value", c.currentSuggestions[c.selectedSuggestion])
			event.Get("target").Call("select")
			c.displayValue = c.currentSuggestions[c.selectedSuggestion]
		}
		env.UnlockRender()

	case "ArrowDown":
		event.Call("preventDefault")
		event.Call("stopPropagation")
		c.selectedSuggestion++
		if c.selectedSuggestion >= len(c.currentSuggestions) {
			c.selectedSuggestion = len(c.currentSuggestions) - 1
			return
		}
		env.Lock()
		if c.selectedSuggestion >= 0 {
			event.Get("target").Set("value", c.currentSuggestions[c.selectedSuggestion])
			event.Get("target").Call("select")
			c.displayValue = c.currentSuggestions[c.selectedSuggestion]
		}
		env.UnlockRender()

	default:
		if c.selectedSuggestion != -1 {
			env.Lock()
			c.selectedSuggestion = -1
			c.currentSuggestions = nil
			env.UnlockRender()
		}

		go func() {
			<-time.After(100 * time.Millisecond)
			target := document.GetElementById(c.editIdStr)
			if target != nil {
				curValue := target.JSValue().Get("value").String()
				c.displayValue = curValue
				if len(curValue) == 0 {
					env.Lock()
					c.currentSuggestions = nil
					env.UnlockRender()
				} else if !c.queryPending && c.Suggestion != nil {
					c.queryPending = true
					c.Suggestion.SuggestionHandle(&suggestionEvent{
						value: curValue,
						env:   env,
						ac:    c,
					})
				}
			}
		}()
	}
}

type suggestionEvent struct {
	value string
	env   vugu.EventEnv
	ac    *AutoComplete
}

func (s *suggestionEvent) Value() string { return s.value }

func (s *suggestionEvent) Propose(strings []string) {
	defer func() {
		s.ac.queryPending = false
	}()
	target := document.GetElementById(s.ac.editIdStr)
	if target != nil {
		curValue := target.JSValue().Get("value").String()
		if curValue != s.value {
			return
		}
		if len(strings) > 50 {
			strings = strings[:50]
		}
		s.env.Lock()
		s.ac.currentSuggestions = strings
		s.env.UnlockRender()
	}
}

type changeEvent struct {
	value string
	env   vugu.EventEnv
	ac    *AutoComplete
}

func (c *changeEvent) Value() string      { return c.value }
func (c *changeEvent) Env() vugu.EventEnv { return c.env }
func (c *changeEvent) SetValue(value string) {
	c.ac.currentSuggestions = nil
	c.ac.selectedSuggestion = -1
	c.ac.Value = value
	c.ac.displayValue = value
	inputElt := document.GetElementById(c.ac.editIdStr)
	if inputElt != nil {
		inputElt.JSValue().Set("value", value)
	}
}

func (c *AutoComplete) Init(vCtx vugu.InitCtx) {
	c.ctx, c.doneFunc = context.WithCancel(context.Background())
	c.displayValue = c.Value
	c.env = vCtx.EventEnv()

	autocompleteIdGen++
	c.editIdStr = fmt.Sprintf("autocomplete-input-%d", autocompleteIdGen)
	c.suggestIdStr = fmt.Sprintf("autocomplete-suggest-%d", autocompleteIdGen)

	c.keyDownFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		c.onKeyDown(args[0], vCtx.EventEnv())
		return nil
	})
}

func (c *AutoComplete) Destroy(vCtx vugu.DestroyCtx) {
	c.doneFunc()
	c.keyDownFunc.Release()
}

type SuggestionEvent interface {
	Value() string

	Propose([]string)
}

type ChangeEvent interface {
	Value() string
	SetValue(string)
	Env() vugu.EventEnv
}

//vugugen:event Suggestion
//vugugen:event Change
