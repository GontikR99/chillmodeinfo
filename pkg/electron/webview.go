// +build wasm,web

package electron

import "github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"

func IsPresent() bool {
	return ipcrenderer.Client != nil
}
