package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/console"
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
)

func init() {
	register(rpcidl.PingHandler(func(message string) {
		console.Log("Ping: "+message)
	}))
}