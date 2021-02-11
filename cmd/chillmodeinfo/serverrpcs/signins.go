// +build server

package serverrpcs

import "github.com/GontikR99/chillmodeinfo/internal/restidl"

func init() {
	register(restidl.HandleLogin())
}