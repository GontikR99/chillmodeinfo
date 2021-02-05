package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
)

func init() {
	register(rpcidl.PingHandler(func(message string) {
		console.Log("Ping: " + message)
	}))
}
