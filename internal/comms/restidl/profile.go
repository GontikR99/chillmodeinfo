package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"net/http"
	"strings"
	"time"
)

const endpointProfileSelfV0 = "/rest/v0/profile"
const endpointProfileAdminV0 = "/rest/v0/profiles"

type Profile interface {
	FetchMine(ctx context.Context) (profile.Entry, error)
	RequestAdmin(ctx context.Context, displayName string) error
	ListAdmins(ctx context.Context) ([]profile.Entry, error)
	UpdateAdmin(ctx context.Context, userId string, state profile.AdminState) error
}

type profileClientStub struct{}

func (p *profileClientStub) FetchMine(ctx context.Context) (profile.Entry, error) {
	req := new(fetchRequestV0)
	res := new(fetchResponseV0)
	err := call(http.MethodGet, endpointProfileSelfV0, req, res)
	return res, err
}

func (p *profileClientStub) RequestAdmin(ctx context.Context, displayName string) error {
	return call(http.MethodPut, endpointProfileSelfV0, &requestAdminRequestV0{displayName}, new(requestAdminResponseV0))
}

func (p *profileClientStub) ListAdmins(ctx context.Context) ([]profile.Entry, error) {
	req := new(listRequestV0)
	res := new(listResponseV0)
	err := call(http.MethodGet, endpointProfileAdminV0, req, res)
	if err != nil {
		return nil, err
	}
	entries := []profile.Entry{}
	for i := 0; i < len(res.Users); i++ {
		entries = append(entries, profile.NewBasicProfile(&res.Users[i]))
	}
	return entries, nil
}

func (p *profileClientStub) UpdateAdmin(ctx context.Context, userId string, newState profile.AdminState) error {
	req := &updateAdminRequestV0{userId, newState}
	res := new(updateAdminResponseV0)
	return call(http.MethodPut, endpointProfileAdminV0, req, res)
}

func GetProfile() Profile {
	return &profileClientStub{}
}

func HandleProfile(handler Profile) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointProfileSelfV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				entry, err := handler.FetchMine(ctx)
				if err == nil {
					return &fetchResponseV0{
						UserId:      entry.GetUserId(),
						Email:       entry.GetEmail(),
						DisplayName: entry.GetDisplayName(),
						AdminState:  entry.GetAdminState(),
						StartDate:   entry.GetStartDate(),
					}, nil
				} else {
					return nil, err
				}
			} else if strings.EqualFold(http.MethodPut, method) {
				req := requestAdminRequestV0{}
				request.ReadTo(&req)
				return &requestAdminResponseV0{}, handler.RequestAdmin(ctx, req.DisplayName)
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
		serve(mux, endpointProfileAdminV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				entries, err := handler.ListAdmins(ctx)
				if err != nil {
					return nil, err
				}
				commEntries := []profile.BasicProfile{}
				for i := 0; i < len(entries); i++ {
					commEntries = append(commEntries, *profile.NewBasicProfile(entries[i]))
				}
				return &listResponseV0{commEntries}, nil
			} else if strings.EqualFold(http.MethodPut, method) {
				req := updateAdminRequestV0{}
				request.ReadTo(&req)
				return &updateAdminResponseV0{}, handler.UpdateAdmin(ctx, req.UserId, req.NewState)
			} else {
				return nil, httputil.UnsupportedMethod(method)
			}
		})
	}
}

type fetchRequestV0 struct{}
type fetchResponseV0 struct {
	UserId      string
	Email       string
	DisplayName string
	AdminState  profile.AdminState
	StartDate   time.Time
}

func (f *fetchResponseV0) GetUserId() string                 { return f.UserId }
func (f *fetchResponseV0) GetStartDate() time.Time           { return f.StartDate }
func (f *fetchResponseV0) GetEmail() string                  { return f.Email }
func (f *fetchResponseV0) GetDisplayName() string            { return f.DisplayName }
func (f *fetchResponseV0) GetAdminState() profile.AdminState { return f.AdminState }

type requestAdminRequestV0 struct{ DisplayName string }
type requestAdminResponseV0 struct{}

type listRequestV0 struct{}
type listResponseV0 struct {
	Users []profile.BasicProfile
}

type updateAdminRequestV0 struct {
	UserId   string
	NewState profile.AdminState
}
type updateAdminResponseV0 struct{}
