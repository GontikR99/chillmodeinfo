// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
)

func init() {
	register(rpcidl.HandleRestartScan(eqfiles.RestartLogScans))
}
