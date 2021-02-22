// +build wasm,electron

package bidoverlay

import (
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/exerpcs"
	"github.com/GontikR99/chillmodeinfo/cmd/electronmain/updateoverlay"
	"github.com/GontikR99/chillmodeinfo/internal/overlay"
	"github.com/GontikR99/chillmodeinfo/internal/profile/localprofile"
)

var currentUpdate *overlay.UpdateEntry

func onOpenBids() {
	currentUpdate=overlay.NewBidUpdate(localprofile.GetProfile(), "","",0)
	currentUpdate.SeqId=updateoverlay.AllocateUpdate()
}

func onCloseBids() {
	if currentUpdate!=nil {
		updateoverlay.Enqueue(currentUpdate.Duplicate())
	}
	currentUpdate=nil
}

func onBid(bidder string, item string, bid float64) {
	if currentUpdate!=nil {
		currentUpdate.Bidder=bidder
		currentUpdate.ItemName=item
		currentUpdate.Bid=bid
	}
}

func init() {
	exerpcs.OnBidOffered(onBid)
}