// +build wasm,web

package roster

import (
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
)

type DumpTarget struct {
	DumpPosted DumpPostedHandler
	Dumps []ParsedDump
}

func (c *DumpTarget) Init(ctx vugu.InitCtx) {
}

func dragOver(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	event.JSEvent().Get("dataTransfer").Set("dropEffect", "copy")
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
			go c.removeDump(event.EventEnv(), dump)
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
