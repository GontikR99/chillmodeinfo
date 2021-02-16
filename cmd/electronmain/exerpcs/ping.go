package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/toast"
	"time"
)

func init() {
	register(rpcidl.PingHandler(func(message string) {
		toast.PopupWithDuration("Ping", message, 10*time.Second)
		console.Log("Ping: " + message)
	}))
}
