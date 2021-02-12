package profile

import "time"

type AdminState int32

const (
	StateAdminUnrequested = AdminState(iota)
	StateAdminRequested
	StateAdminApproved
	StateAdminDenied
	StateAdminMax
)

type Entry interface {
	GetUserId() string
	GetEmail() string
	GetDisplayName() string
	GetAdminState() AdminState
	GetStartDate() time.Time
}

func Equal(oldValue, newValue Entry) bool {
	if oldValue==nil {return newValue==nil}
	if newValue==nil {return oldValue==nil}
	return oldValue.GetUserId()==newValue.GetUserId() &&
		oldValue.GetDisplayName()==newValue.GetDisplayName() &&
		oldValue.GetAdminState()==newValue.GetAdminState()
}

type BasicProfile struct {
	UserId string
	Email string
	DisplayName string
	AdminState AdminState
	StartDate time.Time
}

func (b *BasicProfile) GetUserId() string {return b.UserId}
func (b *BasicProfile) GetEmail() string {return b.Email}
func (b *BasicProfile) GetDisplayName() string {return b.DisplayName}
func (b *BasicProfile) GetAdminState() AdminState {return b.AdminState}
func (b *BasicProfile) GetStartDate() time.Time {return b.StartDate}

