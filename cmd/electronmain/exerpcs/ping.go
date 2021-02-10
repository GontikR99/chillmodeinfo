package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/toast"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"time"
)

func init() {
	register(rpcidl.PingHandler(func(message string) {
		toast.PopupWithDuration("Ping", message, 10*time.Second)
		console.Log("Ping: " + message)
	}))
}
