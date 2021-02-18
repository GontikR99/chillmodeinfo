package record

type Raid interface {
	GetRaidId() uint64
	GetDescription() string
	GetAttendees() []string
	GetDKPValue() float64
}

type BasicRaid struct {
	RaidId      uint64
	Description string
	Attendees   []string
	DKPValue    float64
}

func (b *BasicRaid) GetRaidId() uint64      {return b.RaidId }
func (b *BasicRaid) GetDescription() string {return b.Description}
func (b *BasicRaid) GetAttendees() []string {return b.Attendees}
func (b *BasicRaid) GetDKPValue() float64   {return b.DKPValue}

func NewBasicRaid(evt Raid) *BasicRaid {
	return &BasicRaid{
		RaidId:      evt.GetRaidId(),
		Description: evt.GetDescription(),
		Attendees:   append([]string{}, evt.GetAttendees()...),
		DKPValue:    evt.GetDKPValue(),
	}
}