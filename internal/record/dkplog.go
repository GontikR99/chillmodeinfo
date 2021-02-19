package record

import "time"

type DKPChangeEntry interface {
	// A unique identifier for this entry
	GetEntryId() uint64

	// When did this change occur
	GetTimestamp() time.Time

	// Who had their DKP changed
	GetTarget() string

	// How much of a change
	GetDelta() float64

	// Why changed?
	GetDescription() string

	// Who authorized the change
	GetAuthority() string

	// Raid during which this entry was created
	GetRaidId() uint64
}

type BasicDKPChangeEntry struct {
	EntryId uint64
	Timestamp time.Time
	Target string
	Delta float64
	Description string

	RaidId    uint64
	Authority string
}

func (b *BasicDKPChangeEntry) GetEntryId() uint64 {return b.EntryId}
func (b *BasicDKPChangeEntry) GetTimestamp() time.Time {return b.Timestamp}
func (b *BasicDKPChangeEntry) GetTarget() string {return b.Target}
func (b *BasicDKPChangeEntry) GetDelta() float64 {return b.Delta}
func (b *BasicDKPChangeEntry) GetDescription() string {return b.Description}
func (b *BasicDKPChangeEntry) GetAuthority() string {return b.Authority}
func (b *BasicDKPChangeEntry) GetRaidId() uint64 {return b.RaidId}

func NewBasicDKPChangeEntry(dce DKPChangeEntry) *BasicDKPChangeEntry {
	return &BasicDKPChangeEntry{
		EntryId: dce.GetEntryId(),
		Timestamp:   dce.GetTimestamp(),
		Target:      dce.GetTarget(),
		Delta:       dce.GetDelta(),
		Description: dce.GetDescription(),
		RaidId:      dce.GetRaidId(),
		Authority:   dce.GetAuthority(),
	}
}