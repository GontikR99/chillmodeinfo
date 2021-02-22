// +build wasm,electron

package updateoverlay

import (
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/overlaymap"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
)

var updateIdGen=0

// Reserve an update ID; particularly for bids where the update contents might change
func AllocateUpdate() int {
	updateIdGen++
	return updateIdGen
}

var currentQueue=make(map[int]*overlay.UpdateEntry)

// Add an entry to the update queue
func Enqueue(update *overlay.UpdateEntry) {
	if update.SeqId==0 {
		update.SeqId=AllocateUpdate()
	}
	currentQueue[update.SeqId]=update
}

// Drain the entire update queue
func Drain() map[int]*overlay.UpdateEntry {
	old := currentQueue
	currentQueue = make(map[int]*overlay.UpdateEntry)
	for k, v := range old {
		old[k]=v.Duplicate()
		old[k].Self=profile.NewBasicProfile(localprofile.GetProfile())
	}
	return old
}

type overlayUpdateHandler struct {}

func (o overlayUpdateHandler) Poll() (map[int]*overlay.UpdateEntry, error) {
	return Drain(), nil
}

var uncompletedMap map[int]*overlay.UpdateEntry

func (o overlayUpdateHandler) Enqueue(entries map[int]*overlay.UpdateEntry) error {
	uncompletedMap=entries
	if len(uncompletedMap)==0 {
		om := overlaymap.Lookup("update")
		uw := overlay.Lookup(om.Page)
		if uw!=nil {
			go uw.Close()
		}
	}
	return nil
}
