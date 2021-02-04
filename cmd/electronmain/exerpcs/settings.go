// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
)

func init() {
	register(rpcidl.HandleSetting(settings.LookupSetting, settings.SetSetting))
}