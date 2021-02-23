// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"net/http"
)

func requiresAdmin(ctx context.Context) (profile.Entry, error) {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError != nil {
		return nil, req.IdentityError
	}
	selfProfile, err := dao.LookupProfile(req.UserId)
	if err != nil {
		return nil, err
	}
	if selfProfile.GetAdminState() != profile.StateAdminApproved {
		return nil, httputil.NewError(http.StatusForbidden, "You are not an admin")
	}
	return selfProfile, nil
}
