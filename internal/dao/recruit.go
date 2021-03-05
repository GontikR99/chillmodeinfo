// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

const TableRecruit=db.TableName("recruit")

func TxGetRecruitmentTargets(tx *bbolt.Tx) ([]record.RecruitmentTarget, error) {
	var targets []*recruitTargetV0
	err := db.TxFind(tx, &targets, &bolthold.Query{})
	if err!=nil {
		return nil, err
	}
	var casted []record.RecruitmentTarget
	for _, v := range targets {
		casted=append(casted, v)
	}
	return casted, nil
}

func TxUpsertRecruitmentTarget(tx *bbolt.Tx, r record.RecruitmentTarget) error {
	return db.TxUpsert(tx, r.GetClass(), newRecruitTargetV0(r))
}

type recruitTargetV0 struct {
	Class string `boltholdKey:"Class"`
	Target uint
}

func (r *recruitTargetV0) GetClass() string {return r.Class}
func (r *recruitTargetV0) GetTarget() uint {return r.Target}

func newRecruitTargetV0(r record.RecruitmentTarget) *recruitTargetV0 {
	if r==nil {
		return nil
	} else {
		return &recruitTargetV0{
			Class:  r.GetClass(),
			Target: r.GetTarget(),
		}
	}
}