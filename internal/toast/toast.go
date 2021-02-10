// +build wasm

package toast

import (
	"time"
)

var nextToastId=0
const toastHolderId="toast-holder"

func Error(subsystem string, err error) {
	PopupWithDuration("Error in "+subsystem, err.Error(), 8*time.Second)
}

func Popup(title string, messageText string) {
	PopupWithDuration(title, messageText, time.Duration(0))
}