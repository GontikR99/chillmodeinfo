// +build wasm,web

package roster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/pkg/jsbinding"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"github.com/vugu/vugu"
	"github.com/vugu/vugu/js"
	js2 "syscall/js"
)


func (c *DumpTarget) dropFile(event vugu.DOMEvent) {
	event.StopPropagation()
	event.PreventDefault()
	jsFiles := event.JSEvent().Get("dataTransfer").Get("files")
	if jsFiles.IsUndefined() {
		toast.Error("drop", errors.New("Whoa, was that a file you dropped?"))
	}
	for i:=0;i<jsFiles.Length();i++ {
		go c.handleFile(event.EventEnv(), jsFiles.Index(i))
	}
}

var fileReader=js.Global().Get("FileReader")

func (c *DumpTarget) handleFile(env vugu.EventEnv, fileObj js.Value) {
	fileName := "unspecified"
	if !fileObj.Get("name").IsUndefined() {
		fileName = fileObj.Get("name").String()
	}
	fileSize:=int(-1)
	if !fileObj.Get("size").IsUndefined() {
		fileSize = fileObj.Get("size").Int()
	}
	if fileSize>1024*1024 {
		go c.addDump(env, &uploadError{fileName, fmt.Sprintf("File too large to parse (%2.2f MiB> 1MiB)", float32(fileSize)/(1024*1024))})
		return
	}

	fr := fileReader.New()
	fr.Call("addEventListener", "load", jsbinding.OneShot(func(this js2.Value, args []js2.Value) interface{} {
		e2 := args[0]
		target := e2.Get("target")
		resArr := target.Get("result")
		data := jsbinding.ReadArrayBuffer(resArr)

		members, err := eqfiles.ParseGuildDump(bytes.NewReader(data))
		if err!=nil {
			go c.addDump(env, &uploadError{fileName, err.Error()})
			return nil
		}

		go c.addDump(env, &uploadReady{
			filename: fileName,
			members:  members,
			busy:     false,
		})

		return nil
	}))

	fr.Call("readAsArrayBuffer", fileObj)
}

type uploadError struct {
	filename string
	message string
}

func (u *uploadError) Filename() string {return u.filename}
func (u *uploadError) Message() string {return u.message}
func (u *uploadError) Valid() bool {return false}
func (u *uploadError) Commit(f func(err error)) {}
func (u *uploadError) Busy() bool {return true}

type uploadReady struct {
	filename string
	members []record.Member
	busy bool
}

func (c *uploadReady) Filename() string {return c.filename}

func (c *uploadReady) Message() string {
	mains := 0
	alts := 0
	for _, m := range c.members {
		if m.IsAlt() {
			alts++
		} else {
			mains++
		}
	}
	return fmt.Sprintf("%d mains, %d alts", mains, alts)
}

func (c *uploadReady) Valid() bool {return true}
func (c *uploadReady) Commit(f func(err error)) {
	c.busy=true
	go func() {
		_, err := restidl.Members.MergeMembers(context.Background(), c.members)
		f(err)
	}()
}
func (c *uploadReady) Busy() bool {return c.busy}

