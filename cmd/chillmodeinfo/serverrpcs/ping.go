// build +server

package serverrpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"log"
)

func init() {
	register(restidl.HandlePingV0(func(ping string)string {
		log.Println("Ping: "+ping)
		return "Pong: "+ping
	}))
}