package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"net/http"
	"strings"
	"time"
)

const endpointSelfV0="/rest/v0/profile"

type Profile interface {
	FetchMine(ctx context.Context) (profile.Entry, error)
	RequestAdmin(ctx context.Context, displayName string) error
}

type profileClientStub struct {}

func (p *profileClientStub) FetchMine(ctx context.Context) (profile.Entry, error) {
	req := new(fetchRequestV0)
	res := new(fetchResponseV0)
	err := call(http.MethodGet, endpointSelfV0, req, res)
	return res, err
}

func (p *profileClientStub) RequestAdmin(ctx context.Context, displayName string) error {
	return call(http.MethodPut, endpointSelfV0, &requestAdminRequestV0{displayName}, new(requestAdminResponseV0))
}


func GetProfile() Profile {
	return &profileClientStub{}
}


func HandleProfile(handler Profile) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, endpointSelfV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold(http.MethodGet, method) {
				entry, err := handler.FetchMine(ctx)
				if err==nil {
					return &fetchResponseV0{
						IdToken:     entry.GetUserId(),
						Email:       entry.GetEmail(),
						DisplayName: entry.GetDisplayName(),
						AdminState:  entry.GetAdminState(),
						StartDate: entry.GetAdminStartDate(),
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
	}
}

type fetchRequestV0 struct {}
type fetchResponseV0 struct {
	IdToken     string
	Email       string
	DisplayName string
	AdminState  profile.AdminState
	StartDate   time.Time
}

func (f *fetchResponseV0) GetUserId() string {return f.IdToken}
func (f *fetchResponseV0) GetAdminStartDate() time.Time {return f.StartDate}
func (f *fetchResponseV0) GetEmail() string               {return f.Email}
func (f *fetchResponseV0) GetDisplayName() string            {return f.DisplayName}
func (f *fetchResponseV0) GetAdminState() profile.AdminState {return f.AdminState}

type requestAdminRequestV0 struct {DisplayName string}
type requestAdminResponseV0 struct {}
