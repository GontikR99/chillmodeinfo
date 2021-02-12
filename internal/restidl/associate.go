package restidl

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"net/http"
	"strings"
)

const pathAssociateV0 = "/rest/v0/associate"

type AssociationHandler interface {
	AssertAssociatedClientId(clientId string) error
	AssociateClientId(clientId string, req *Request) error
	DisassociateClientId(req *Request) error
}

type assertAssociatedRequestV0 struct{ ClientId string }
type assertAssociatedResponseV0 struct{}

// Check to see if there exists an association for the given clientId, returning an error if there is not
func AssertAssociated(clientId string) error {
	req := &assertAssociatedRequestV0{ClientId: clientId}
	res := new(assertAssociatedResponseV0)
	return call("GET", pathAssociateV0, req, res)
}

type associateRequestV0 struct{ ClientId string }
type associateResponseV0 struct{}

// Assuming the user is logged in, associate the specified clientId with the user
func AssociateClientId(clientId string) error {
	req := &associateRequestV0{ClientId: clientId}
	res := new(associateResponseV0)
	return call("PUT", pathAssociateV0, req, res)
}

type disassociateRequestV0 struct{}
type disassociateResponseV0 struct{}

// Assuming the user is logged in with a clientId, delete the association between the clientId and identity
func DisassociateClientId() error {
	req := &disassociateRequestV0{}
	res := new(disassociateResponseV0)
	return call("DELETE", pathAssociateV0, req, res)
}

func HandleAssociate(handler AssociationHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		serve(mux, pathAssociateV0, func(ctx context.Context, method string, request *Request) (interface{}, error) {
			if strings.EqualFold("GET", method) {
				reqVal := new(assertAssociatedRequestV0)
				request.ReadTo(reqVal)
				err := handler.AssertAssociatedClientId(reqVal.ClientId)
				return new(assertAssociatedResponseV0), err
			} else if strings.EqualFold("PUT", method) {
				reqVal := new(associateRequestV0)
				request.ReadTo(reqVal)
				err := handler.AssociateClientId(reqVal.ClientId, request)
				return new(associateResponseV0), err
			} else if strings.EqualFold("DELETE", method) {
				reqVal := new(disassociateRequestV0)
				request.ReadTo(reqVal)
				err := handler.DisassociateClientId(request)
				return new(disassociateResponseV0), err

			} else {
				return nil, httputil.NewError(http.StatusBadRequest, "Bad method for call to associate")
			}
		})
	}
}
