package record

import (
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"time"
)

type Member interface {
	GetName() string
	GetClass() string
	GetLevel() int16
	GetRank() string
	IsAlt() bool
	GetDKP() float64
	GetLastActive() time.Time
	GetOwner() string
}

func IsActive(m Member) bool {
	return !m.GetLastActive().Before(time.Now().Add(sitedef.InactiveDuration))
}

type BasicMember struct {
	Name       string
	Class      string
	Level      int16
	Rank       string
	Alt        bool
	DKP        float64
	LastActive time.Time
	Owner      string
}

func (b *BasicMember) GetName() string          { return b.Name }
func (b *BasicMember) GetClass() string         { return b.Class }
func (b *BasicMember) GetLevel() int16          { return b.Level }
func (b *BasicMember) GetRank() string          { return b.Rank }
func (b *BasicMember) IsAlt() bool              { return b.Alt }
func (b *BasicMember) GetDKP() float64          { return b.DKP }
func (b *BasicMember) GetLastActive() time.Time { return b.LastActive }
func (b *BasicMember) GetOwner() string         { return b.Owner }

func NewBasicMember(member Member) *BasicMember {
	if member == nil {
		return nil
	}
	return &BasicMember{
		Name:       member.GetName(),
		Class:      member.GetClass(),
		Level:      member.GetLevel(),
		Rank:       member.GetRank(),
		Alt:        member.IsAlt(),
		DKP:        member.GetDKP(),
		LastActive: member.GetLastActive(),
		Owner:      member.GetOwner(),
	}
}
