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
	GetAdminStartDate() time.Time
}

func Equal(oldValue, newValue Entry) bool {
	if oldValue==nil {return newValue==nil}
	if newValue==nil {return oldValue==nil}
	return oldValue.GetUserId()==newValue.GetUserId() &&
		oldValue.GetDisplayName()==newValue.GetDisplayName() &&
		oldValue.GetAdminState()==newValue.GetAdminState()
}