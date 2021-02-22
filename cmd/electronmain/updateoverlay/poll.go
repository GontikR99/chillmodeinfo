// +build wasm,electron

package updateoverlay

import (
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/overlaymap"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"time"
)

func PollForUpdates() {
	go func() {
		for {
			<-time.After(100*time.Millisecond)
			if len(currentQueue)==0 {
				continue
			}
			om := overlaymap.Lookup("update")
			uw := overlay.Lookup(om.Page)
			if uw==nil {
				bw := overlay.Launch(om.Page, true)
				server := exerpcs.NewServer()
				bw.ServeRPC(server)

				bw.OnClosed(func() {
					if uncompletedMap!=nil {
						for k, v := range uncompletedMap {
							v.SeqId = k
							Enqueue(v)
						}
						uncompletedMap=nil
					}
				})

				bw.JSValue().Get("webContents").Call("openDevTools", map[string]interface{} {
					"mode":"detach",
					"activate":"false",
				})
			}
		}
	}()
}