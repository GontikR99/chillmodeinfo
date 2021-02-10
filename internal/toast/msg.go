// +build wasm

package toast

import "time"

type toastMessage struct {
	Title string
	Body string
	Timeout time.Duration
}

const channelToast = "toasts"

