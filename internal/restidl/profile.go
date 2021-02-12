package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"net/http"
	"strings"
)

const endpointSelfV0="/rest/v0/profile"

type Profile interface {
	FetchMine(ctx context.Context) (profile.Entry, error)
}

type profileClientStub struct {}

func (p *profileClientStub) FetchMine(ctx context.Context) (profile.Entry, error) {
	req := new(fetchRequestV0)
	res := new(fetchResponseV0)
	err := call(http.MethodGet, endpointSelfV0, req, res)
	return res, err
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
						IdToken:     entry.GetIdToken(),
						Email:       entry.GetEmail(),
						DisplayName: entry.GetDisplayName(),
						AdminState:  entry.GetAdminState(),
					}, nil
				} else {
					return nil, err
				}
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
}

func (f *fetchResponseV0) GetIdToken() string  {return f.IdToken }
func (f *fetchResponseV0) GetEmail() string               {return f.Email }
func (f *fetchResponseV0) GetDisplayName() string            {return f.DisplayName }
func (f *fetchResponseV0) GetAdminState() profile.AdminState {return f.AdminState }
