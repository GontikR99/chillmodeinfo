// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"go.etcd.io/bbolt"
	"net/http"
	"time"
)

type serverDKPLogHandler struct {}

func (s serverDKPLogHandler) Append(ctx context.Context, delta record.DKPChangeEntry) error {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError != nil {
		return req.IdentityError
	}
	selfProfile, err := dao.LookupProfile(req.UserId)
	if err != nil {
		return err
	}
	if selfProfile.GetAdminState()!=profile.StateAdminApproved {
		return httputil.NewError(http.StatusForbidden, "Only admins may append log entries")
	}

	return db.MakeUpdate(func(tx *bbolt.Tx) error {
		target:=initialCap(delta.GetTarget())
		targetMemberRecord, err := dao.GetMember(target)
		if err!=nil {
			return err
		}
		nowTime:=time.Now()
		updatedMember :=record.NewBasicMember(targetMemberRecord)
		updatedMember.LastActive=nowTime
		updatedMember.DKP = updatedMember.DKP+ delta.GetDelta()

		updatedDelta:=record.NewBasicDKPChangeEntry(delta)
		updatedDelta.Target=updatedMember.GetName()
		updatedDelta.Authority=selfProfile.GetDisplayName()
		updatedDelta.Timestamp=nowTime

		err=dao.TxUpsertMember(tx, updatedMember)
		if err!=nil {
			return err
		}
		return dao.TxAppendDKPChange(tx, updatedDelta)
	})
}

func (s serverDKPLogHandler) Retrieve(ctx context.Context, target string) ([]record.DKPChangeEntry, error) {
	return dao.GetDKPChangesForTarget(target)
}

func (s serverDKPLogHandler) Remove(ctx context.Context, entryId uint64) error {
	req := ctx.Value(restidl.TagRequest).(*restidl.Request)
	if req.IdentityError != nil {
		return req.IdentityError
	}
	selfProfile, err := dao.LookupProfile(req.UserId)
	if err != nil {
		return err
	}
	if selfProfile.GetAdminState()!=profile.StateAdminApproved {
		return httputil.NewError(http.StatusForbidden, "Only admins may append log entries")
	}

	return db.MakeUpdate(func(tx *bbolt.Tx) error {
		dkpEntry, err := dao.TxGetDKPChange(tx, entryId)
		if err!=nil {
			return err
		}
		if dkpEntry.GetRaidId()!=0 {
			return httputil.NewError(http.StatusBadRequest, "You may not remove a log entry corresponding to a raid.  Go remove the raid.")
		}
		dao.TxRemoveDKPChange(tx, entryId)

		member, err := dao.GetMember(dkpEntry.GetTarget())
		if err!=nil {
			return err
		}

		newMember := record.NewBasicMember(member)
		newMember.DKP -= dkpEntry.GetDelta()
		newMember.LastActive = time.Now()

		return dao.TxUpsertMember(tx, newMember)
	})
}

func init() {
	register(restidl.HandleDKPLog(&serverDKPLogHandler{}))
}