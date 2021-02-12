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

type serverProfileHandler struct {}


func (s *serverProfileHandler) UpdateAdmin(ctx context.Context, userId string, state profile.AdminState) error {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError!=nil {
		return req.IdentityError
	}
	selfProfile, err := dao.LookupProfile(req.UserId)
	if err!=nil {
		return err
	}
	if selfProfile.GetAdminState()!=profile.StateAdminApproved {
		return httputil.NewError(http.StatusForbidden, "You are not an admin")
	}
	userProfile, err := dao.LookupProfile(userId)
	if err!=nil {
		return err
	}
	if userProfile.GetAdminState()==profile.StateAdminUnrequested {
		return httputil.NewError(http.StatusBadRequest, "This user did not request promotion")
	}
	if userProfile.GetAdminState()==profile.StateAdminApproved && selfProfile.GetStartDate().After(userProfile.GetStartDate()) {
		return httputil.NewError(http.StatusForbidden, "You may not modify the state of admins elder to yourself")
	}
	dao.UpdateProfileForAdmin(userId, userProfile.GetDisplayName(), state, time.Now())
	return nil
}

func (s *serverProfileHandler) ListAdmins(ctx context.Context) ([]profile.Entry, error) {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError!=nil {
		return nil, req.IdentityError
	}
	selfProfile, err := dao.LookupProfile(req.UserId)
	if err!=nil {
		return nil, err
	}
	if selfProfile.GetAdminState()!=profile.StateAdminApproved {
		return nil, httputil.NewError(http.StatusForbidden, "You are not an admin")
	}

	resultProfiles:=[]profile.Entry{}
	allProfiles := dao.ListAllProfiles()
	for i:=0;i<len(allProfiles);i++ {
		if allProfiles[i].GetAdminState()==profile.StateAdminUnrequested {
			continue
		}
		resultProfile:=&profile.BasicProfile{
			UserId:      allProfiles[i].GetUserId(),
			DisplayName: allProfiles[i].GetDisplayName(),
			AdminState:  allProfiles[i].GetAdminState(),
			StartDate:   allProfiles[i].GetStartDate(),
		}
		if allProfiles[i].GetAdminState()==profile.StateAdminRequested || allProfiles[i].GetAdminState()>=profile.StateAdminDenied {
			resultProfile.Email=allProfiles[i].GetEmail()
		} else {
			resultProfile.Email="undisclosed"
		}
		resultProfiles=append(resultProfiles, resultProfile)
	}
	return resultProfiles, nil
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