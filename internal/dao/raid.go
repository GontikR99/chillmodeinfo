// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"time"
)

func TxInsertRaid(tx *bbolt.Tx, raid record.Raid) (uint64, error) {
	raidRec := newRaidV0(raid)
	raidRec.RaidId=0

	err := db.TxInsert(tx, bolthold.NextSequence(), raidRec)
	return raidRec.RaidId, err
}

func TxGetRaids(tx *bbolt.Tx) ([]record.Raid, error) {
	raidResult:=new([]record.Raid)
	err := db.TxForEach(tx, &bolthold.Query{}, func(raid *raidV0)error {
		*raidResult=append(*raidResult, raid)
		return nil
	})
	return *raidResult, err
}

func TxGetRaid(tx *bbolt.Tx, raidId uint64) (record.Raid, error) {
	var raid raidV0
	err := db.TxGet(tx, raidId, &raid)
	return &raid, err
}

func TxDeleteRaid(tx *bbolt.Tx, raidId uint64) error {
	return db.TxDelete(tx, raidId, &raidV0{})
}

type raidV0 struct {
	RaidId      uint64 `boltholdKey:"RaidId"`
	Timestamp time.Time
	Description string
	Attendees   []string
	DKPValue    float64
}

func (b *raidV0) GetRaidId() uint64      {return b.RaidId }
func (b *raidV0) GetTimestamp() time.Time {return b.Timestamp}
func (b *raidV0) GetDescription() string {return b.Description}
func (b *raidV0) GetAttendees() []string {return b.Attendees}
func (b *raidV0) GetDKPValue() float64   {return b.DKPValue}

func newRaidV0(evt record.Raid) *raidV0 {
	return &raidV0{
		RaidId:      evt.GetRaidId(),
		Timestamp: evt.GetTimestamp(),
		Description: evt.GetDescription(),
		Attendees:   append([]string{}, evt.GetAttendees()...),
		DKPValue:    evt.GetDKPValue(),
	}
}