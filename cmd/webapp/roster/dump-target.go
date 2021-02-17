// +build wasm,web

package roster

import (
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"syscall/js"
	"time"
)

type DumpTarget struct {
	DumpPosted DumpPostedHandler
	Dumps []ParsedDump
	dragOverFunc js.Func
	dead bool
}

func (c *DumpTarget) Init(ctx vugu.InitCtx) {
	c.dragOverFunc=js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		event.Call("stopPropagation")
		event.Call("preventDefault")
		event.Get("dataTransfer").Set("dropEffect", "copy")
		return nil
	})
	c.dead=false
	go func() {
		for !c.dead {
			elt := document.GetElementById("guild-dump-drop")
			if elt!=nil {
				elt.AddEventListener("dragover", c.dragOverFunc)
				return
			}
			<-time.After(10*time.Millisecond)
		}
	}()
}

func (c *DumpTarget) Destroy(ctx vugu.DestroyCtx) {
	c.dead = true
	c.dragOverFunc.Release()
}

func (c *DumpTarget) addDump(env vugu.EventEnv, dump ParsedDump) {
	newDumps := []ParsedDump{dump}
	newDumps = append(newDumps, c.Dumps...)
	env.Lock()
	c.Dumps=newDumps
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
	Filename() string
	Message() string
	Valid() bool
	Commit(func(err error))
	Busy()bool
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
