// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"go.etcd.io/bbolt"
	"net/http"
)

type serverRecruitHandler struct {}

func (s serverRecruitHandler) Update(ctx context.Context, target record.RecruitmentTarget) error {
	_, err := requiresAdmin(ctx)
	if err!=nil {return err}
	if _, isClazz := eqspec.ClassMap[target.GetClass()]; !isClazz {
		return httputil.NewError(http.StatusBadRequest, target.GetClass()+" is not a recognized class")
	}
	return db.MakeUpdate([]db.TableName{dao.TableRecruit}, func(tx *bbolt.Tx) error {
		return dao.TxUpsertRecruitmentTarget(tx, target)
	})
}

func (s serverRecruitHandler) Fetch(ctx context.Context) ([]record.RecruitmentTarget, error) {
	var targets []record.RecruitmentTarget
	err := db.MakeView([]db.TableName{dao.TableRecruit}, func(tx *bbolt.Tx) error {
		var err error
		targets, err = dao.TxGetRecruitmentTargets(tx)
		return err
	})
	return targets, err
}

func init() {
	register(restidl.HandleRecruit(serverRecruitHandler{}))
}