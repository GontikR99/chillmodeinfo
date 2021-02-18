// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

func TxAddRaid(tx *bbolt.Tx, raid record.Raid) (uint64, error) {
	raidRec := newRaidV0(raid)
	raidRec.RaidId=0

	err := db.Insert(bolthold.NextSequence(), raidRec)
	return raidRec.RaidId, err
}


type raidV0 struct {
	RaidId      uint64 `boltholdKey:"RaidId"`
	Description string
	Attendees   []string
	DKPValue    float64
}

func (b *raidV0) GetRaidId() uint64      {return b.RaidId }
func (b *raidV0) GetDescription() string {return b.Description}
func (b *raidV0) GetAttendees() []string {return b.Attendees}
func (b *raidV0) GetDKPValue() float64   {return b.DKPValue}

func newRaidV0(evt record.Raid) *raidV0 {
	return &raidV0{
		RaidId:      evt.GetRaidId(),
		Description: evt.GetDescription(),
		Attendees:   append([]string{}, evt.GetAttendees()...),
		DKPValue:    evt.GetDKPValue(),
	}
}