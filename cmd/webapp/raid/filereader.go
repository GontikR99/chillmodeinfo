// +build wasm,web

package raid

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
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
	for i := 0; i < jsFiles.Length(); i++ {
		c.handleFile(event.EventEnv(), jsFiles.Index(i))
	}
}

var fileReader = js.Global().Get("FileReader")

func (c *DumpTarget) handleFile(env vugu.EventEnv, fileObj js.Value) {
	fileName := "unspecified"
	if !fileObj.Get("name").IsUndefined() {
		fileName = fileObj.Get("name").String()
	}
	fileSize := int(-1)
	if !fileObj.Get("size").IsUndefined() {
		fileSize = fileObj.Get("size").Int()
	}
	if fileSize > 1024*1024 {
		go c.addDump(env, &uploadError{
			uid:      newUID(),
			filename: fileName,
			message:  fmt.Sprintf("File too large to parse (%2.2f MiB> 1MiB)", float32(fileSize)/(1024*1024))})
		return
	}

	fr := fileReader.New()
	fr.Call("addEventListener", "load", jsbinding.OneShot(func(this js2.Value, args []js2.Value) interface{} {
		e2 := args[0]
		target := e2.Get("target")
		resArr := target.Get("result")
		data := jsbinding.ReadArrayBuffer(resArr)

		attendees, err := eqspec.ParseRaidDump(bytes.NewReader(data))
		if err != nil {
			go c.addDump(env, &uploadError{
				uid:      newUID(),
				filename: fileName,
				message:  err.Error()})
			return nil
		}
		go c.addDump(env, &uploadReady{
			uid:       newUID(),
			filename:  fileName,
			attendees: attendees,
			busy:      false,
		})
		return nil
	}))

	fr.Call("readAsArrayBuffer", fileObj)
}

type raidInfoer struct {
	description string
	dkp         float64
}

func (r *raidInfoer) Description() string     { return r.description }
func (r *raidInfoer) SetDescription(d string) { r.description = d }
func (r *raidInfoer) DKP() float64            { return r.dkp }
func (r *raidInfoer) SetDKP(d float64)        { r.dkp = d }

type uploadError struct {
	raidInfoer
	uid      string
	filename string
	message  string
}

func (u *uploadError) UniqueId() string         { return u.uid }
func (u *uploadError) Filename() string         { return u.filename }
func (u *uploadError) Message() string          { return u.message }
func (u *uploadError) Valid() bool              { return false }
func (u *uploadError) Commit(f func(err error)) {}
func (u *uploadError) Busy() bool               { return true }

type uploadReady struct {
	raidInfoer
	uid       string
	filename  string
	attendees []string
	busy      bool
}

func (c *uploadReady) UniqueId() string { return c.uid }
func (c *uploadReady) Filename() string { return c.filename }

func (c *uploadReady) Message() string {
	return fmt.Sprintf("%d attendees", len(c.attendees))
}

func (c *uploadReady) Valid() bool {
	return c.description != "" && c.dkp > 0
}
func (c *uploadReady) Commit(f func(err error)) {
	c.busy = true
	go func() {
		err := restidl.Raid.Add(context.Background(), &record.BasicRaid{
			Description: c.description,
			Attendees:   c.attendees,
			DKPValue:    c.dkp,
		})
		f(err)
	}()
}
func (c *uploadReady) Busy() bool { return c.busy }

func newUID() string {
	bb := make([]byte, 20)
	rand.Read(bb)
	return hex.EncodeToString(bb)
}
