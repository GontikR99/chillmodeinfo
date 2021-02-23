package record

import "time"

type Raid interface {
	GetRaidId() uint64
	GetTimestamp() time.Time
	GetDescription() string
	GetAttendees() []string
	GetDKPValue() float64
	GetCredited() []string
}

type BasicRaid struct {
	RaidId      uint64
	Timestamp   time.Time
	Description string
	Attendees   []string
	DKPValue    float64
	Credited    []string
}

func (b *BasicRaid) GetRaidId() uint64       { return b.RaidId }
func (b *BasicRaid) GetTimestamp() time.Time { return b.Timestamp }
func (b *BasicRaid) GetDescription() string  { return b.Description }
func (b *BasicRaid) GetAttendees() []string  { return b.Attendees }
func (b *BasicRaid) GetDKPValue() float64    { return b.DKPValue }
func (b *BasicRaid) GetCredited() []string   { return b.Credited }

func NewBasicRaid(evt Raid) *BasicRaid {
	if evt == nil {
		return nil
	}
	return &BasicRaid{
		RaidId:      evt.GetRaidId(),
		Timestamp:   evt.GetTimestamp(),
		Description: evt.GetDescription(),
		Attendees:   append([]string{}, evt.GetAttendees()...),
		DKPValue:    evt.GetDKPValue(),
		Credited:    append([]string{}, evt.GetCredited()...),
	}
}
