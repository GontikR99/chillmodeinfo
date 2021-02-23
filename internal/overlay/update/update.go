// +build wasm

package update

import (
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
)

type UpdateType int

const (
	UpdateGuildDump = UpdateType(iota)
	UpdateRaidDump
	UpdateBid
	UpdateError
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
	Attendees []string

	// For loot
	Bidder string
	ItemName string
	Bid float64

	// For error message
	ErrorMessage string
}

func NewGuildDump(members []record.Member) *UpdateEntry {
	var newMembers []*record.BasicMember
	for _, v := range members {
		newMembers = append(newMembers, record.NewBasicMember(v))
	}
	return &UpdateEntry{
		Type:    UpdateGuildDump,
		SeqId:   0,
		Members: newMembers,
	}
}

func NewRaidDump(attendees []string) *UpdateEntry {
	return &UpdateEntry{
		Type:      UpdateRaidDump,
		SeqId:     0,
		Attendees: attendees,
	}
}

func NewBidUpdate(bidder string, itemName string, bid float64) *UpdateEntry {
	return &UpdateEntry{
		Type:     UpdateBid,
		SeqId:    0,
		Bidder:   bidder,
		ItemName: itemName,
		Bid:      bid,
	}
}

func NewError(message string) *UpdateEntry {
	return &UpdateEntry{
		Type:         UpdateError,
		SeqId:        0,
		ErrorMessage: message,
	}
}

func (u *UpdateEntry) Duplicate() *UpdateEntry {
	return &UpdateEntry{
		Type:         u.Type,
		SeqId:        u.SeqId,
		Self:         u.Self,
		Members:      u.Members,
		Attendees:    u.Attendees,
		Bidder:       u.Bidder,
		ItemName:     u.ItemName,
		Bid:          u.Bid,
		ErrorMessage: u.ErrorMessage,
	}
}