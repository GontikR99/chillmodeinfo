// +build server

package serverrpcs

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile/signins"
	"net/http"
	"strings"
)

type serverAssociateHandler struct{}

func (s *serverAssociateHandler) AssertAssociatedClientId(clientId string) error {
	_, present, err := dao.LookupClientId(clientId)
	if err != nil {
		return err
	} else if !present {
		return httputil.NewError(http.StatusNotFound, "No such clientId associated")
	} else {
		return nil
	}
}

func (s *serverAssociateHandler) AssociateClientId(clientId string, req *restidl.Request) error {
	if req.IdentityError != nil {
		return req.IdentityError
	} else if req.UserId == "" {
		return errors.New("No userId to associate")
	} else {
		return dao.AssociateClientId(clientId, req.UserId)
	}
}

func (s *serverAssociateHandler) DisassociateClientId(req *restidl.Request) error {
	if req.IdentityError != nil {
		return req.IdentityError
	} else if !strings.HasPrefix(req.IdToken, signins.TokenClientId) {
		return httputil.NewError(http.StatusBadRequest, "Can only disassociate clientIds")
	} else {
		return dao.DisassociateClientId(req.IdToken[len(signins.TokenClientId):])
	}
}

func init() {
	register(restidl.HandleAssociate(&serverAssociateHandler{}))
}
