// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
)

func init() {
	register(rpcidl.HandleRestartScan(eqspec.RestartLogScans))
}
