// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
	"net/http"
	"time"
)

type serverProfileHandler struct {
}

func (s *serverProfileHandler) FetchMine(ctx context.Context) (profile.Entry, error) {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError!=nil {
		return nil, req.IdentityError
	} else {
		return dao.LookupProfile(req.UserId)
	}
}

func (s *serverProfileHandler) RequestAdmin(ctx context.Context, displayName string) error {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError!=nil {
		return req.IdentityError
	}
	entry, err := dao.LookupProfile(req.UserId)
	if err!=nil {
		return err
	}
	if entry.GetAdminState()!=profile.StateAdminUnrequested {
		return httputil.NewError(http.StatusForbidden, "You may not request admin privileges at this time")
	}
	if displayName=="" {
		return httputil.NewError(http.StatusBadRequest, "You must have a non-empty name to become an admin")
	}
	return dao.UpdateProfileForAdmin(req.UserId, displayName, profile.StateAdminRequested, time.Now())
}


func init() {
	register(restidl.HandleProfile(&serverProfileHandler{}))
}