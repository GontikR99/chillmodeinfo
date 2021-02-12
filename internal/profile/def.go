package profile

type AdminState int32

const (
	StateAdminUnrequested = AdminState(iota)
	StateAdminRequested
	StateAdminApproved
	StateAdminDenied
	StateAdminMax
)

type Entry interface {
	GetIdToken() string
	GetEmail() string
	GetDisplayName() string
	GetAdminState() AdminState
}

func Equal(oldValue, newValue Entry) bool {
	if oldValue==nil {return newValue==nil}
	if newValue==nil {return oldValue==nil}
	return oldValue.GetIdToken()==newValue.GetIdToken() &&
		oldValue.GetDisplayName()==newValue.GetDisplayName() &&
		oldValue.GetAdminState()==newValue.GetAdminState()
}