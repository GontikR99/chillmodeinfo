// +build wasm,web

package raid

import (
	"context"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"sort"
	"strconv"
	"syscall/js"
	"time"
)

type DumpTarget struct {
	DumpPosted DumpPostedHandler
	Dumps []ParsedDump
	dragOverFunc js.Func

	ctx context.Context
	ctxDone context.CancelFunc
}

func (c *DumpTarget) Init(vCtx vugu.InitCtx) {
	c.ctx, c.ctxDone = context.WithCancel(context.Background())

	c.dragOverFunc=js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		event.Call("stopPropagation")
		event.Call("preventDefault")
		event.Get("dataTransfer").Set("dropEffect", "copy")
		return nil
	})
	go func() {
		for {
			elt := document.GetElementById("raid-dump-drop")
			if elt!=nil {
				elt.AddEventListener("dragover", c.dragOverFunc)
				return
			}
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(10 * time.Millisecond):
			}
		}
	}()
}

func (c *DumpTarget) Destroy(vCtx vugu.DestroyCtx) {
	c.ctxDone()
	c.dragOverFunc.Release()
}

type byFilename []ParsedDump
func (b byFilename) Len() int {return len(b)}
func (b byFilename) Less(i, j int) bool {return b[i].Filename()<b[j].Filename()}
func (b byFilename) Swap(i, j int) {b[i],b[j] = b[j],b[i]}

func (c *DumpTarget) addDump(env vugu.EventEnv, dump ParsedDump) {
	env.Lock()
	c.Dumps=append(c.Dumps, dump)
	sort.Stable(byFilename(c.Dumps))
	env.UnlockRender()
}

func (c *DumpTarget) removeDump(env vugu.EventEnv, dump ParsedDump) {
	var newDumps []ParsedDump
	for _, v := range c.Dumps {
		if v!=dump {
			newDumps=append(newDumps, v)
		}
	}
	env.Lock()
	c.Dumps=newDumps
	env.UnlockRender()
}

func (c *DumpTarget) Commit(event vugu.DOMEvent, dump ParsedDump) {
	event.PreventDefault()
	dump.Commit(func(err error) {
		if err==nil {
			go func() {
				c.removeDump(event.EventEnv(), dump)
				c.DumpPosted.DumpPostedHandle(DumpPostedEvent{Env: event.EventEnv()})
			}()
		} else {
			toast.Error("uploads", err)
		}
	})
}

func (c *DumpTarget) Abort(event vugu.DOMEvent, dump ParsedDump) {
	event.PreventDefault()
	go c.removeDump(event.EventEnv(), dump)
}

type ParsedDump interface {
	UniqueId() string
	Filename() string
	Message() string
	Valid() bool
	Busy()bool

	Description()string
	SetDescription(string)

	DKP()float64
	SetDKP(float64)

	Commit(donefunc func(err error))
}

func (c *DumpTarget) textFocus(event vugu.DOMEvent) {
	event.JSEventTarget().Call("select")
}

func (c *DumpTarget) descriptionChange(event vugu.DOMEvent, dump ParsedDump) {
	event.PreventDefault()
	dump.SetDescription(event.JSEventTarget().Get("value").String())
}

func (c *DumpTarget) dkpChange(event vugu.DOMEvent, dump ParsedDump) {
	event.PreventDefault()
	dkpTxt := event.JSEventTarget().Get("value").String()
	dkp, err := strconv.ParseFloat(dkpTxt, 64)
	if err!=nil {
		toast.Error("DKP update", errors.New("Please enter a number"))
		event.JSEventTarget().Set("value", fmt.Sprintf("%.1f", dkp))
	} else {
		dump.SetDKP(dkp)
	}
}

type dumpAttrs struct {
	dump ParsedDump
}

func (d dumpAttrs) AttributeList() []vugu.VGAttribute {
	if d.dump.Busy() {
		return []vugu.VGAttribute{{"","disabled", "true"}}
	} else {
		return nil
	}
}

type DumpPostedEvent struct {
	Env vugu.EventEnv
}

//vugugen:event DumpPosted
