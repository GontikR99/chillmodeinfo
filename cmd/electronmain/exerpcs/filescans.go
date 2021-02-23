// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
)

func init() {
	register(rpcidl.HandleRestartScan(eqspec.RestartLogScans))
}
