// +build wasm

package overlay

import (
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
)

type UpdateType int

const (
	UpdateGuildDump = UpdateType(iota)
	UpdateRaidDump
	UpdateBid
)

type UpdateEntry struct {
	// Type of the update
	Type UpdateType

	// Sequence number of the update
	SeqId int

	// Who I am, just in case that hasn't propagated yet.
	Self *profile.BasicProfile

	// For guild dumps
	Members []*record.BasicMember

	// For raid dumps
	Raid *record.BasicRaid

	// For loot
	Bidder string
	ItemName string
	Bid float64
}

func NewGuildDump(entry profile.Entry, members []record.Member) *UpdateEntry {
	var newMembers []*record.BasicMember
	for _, v := range members {
		newMembers = append(newMembers, record.NewBasicMember(v))
	}
	return &UpdateEntry{
		Type:    UpdateGuildDump,
		SeqId:   0,
		Self:    profile.NewBasicProfile(entry),
		Members: newMembers,
	}
}

func NewRaidDump(entry profile.Entry, raid record.Raid) *UpdateEntry {
	return &UpdateEntry{
		Type:  UpdateRaidDump,
		SeqId: 0,
		Self:  profile.NewBasicProfile(entry),
		Raid:  record.NewBasicRaid(raid),
	}
}

func NewBidUpdate(entry profile.Entry, bidder string, itemName string, bid float64) *UpdateEntry {
	return &UpdateEntry{
		Type:     UpdateBid,
		SeqId:    0,
		Self:     profile.NewBasicProfile(entry),
		Bidder:   bidder,
		ItemName: itemName,
		Bid:      bid,
	}
}

func (u *UpdateEntry) Duplicate() *UpdateEntry {
	return &UpdateEntry{
		Type:     u.Type,
		SeqId:    u.SeqId,
		Self:     u.Self,
		Members:  u.Members,
		Raid:     u.Raid,
		Bidder:   u.Bidder,
		ItemName: u.ItemName,
		Bid:      u.Bid,
	}
}