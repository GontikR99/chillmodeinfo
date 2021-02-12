// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/restidl"
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

func init() {
	register(restidl.HandleProfile(&serverProfileHandler{}))
}