// +build wasm,web

package util

import (
	"fmt"
	"syscall/js"
	"time"
)


var startDate=js.Global().Get("Date").New()
var tzOffset=startDate.Call("getTimezoneOffset").Int()

func FormatTimestamp(timestamp time.Time) string {
	nowTime := time.Now()
	if timestamp.After(nowTime) {
		duration:=timestamp.Sub(nowTime)
		return fmt.Sprintf("Pending for %dh%02dm", int(duration.Hours()), int(duration.Minutes())%60)
	}
	return timestamp.Add(-time.Duration(tzOffset)*time.Minute).Format("2006-01-02 03:04 PM")
}