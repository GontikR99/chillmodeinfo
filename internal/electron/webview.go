// +build wasm,web

package electron

import "github.com/GontikR99/chillmodeinfo/internal/electron/ipc/ipcrenderer"

func IsPresent() bool {
	return ipcrenderer.Client!=nil
}