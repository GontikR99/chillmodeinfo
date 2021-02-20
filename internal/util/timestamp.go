// +build wasm,web

package util

import (
	"syscall/js"
	"time"
)


var startDate=js.Global().Get("Date").New()
var tzOffset=startDate.Call("getTimezoneOffset").Int()

func FormatTimestamp(timestamp time.Time) string {
	return timestamp.Add(-time.Duration(tzOffset)*time.Minute).Format("2006-01-02 03:04 PM")
}