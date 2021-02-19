// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"sort"
	"time"
)

func TxAppendDKPChange(tx *bbolt.Tx, dcle record.DKPChangeEntry) error {
	return db.TxInsert(tx, bolthold.NextSequence(), newDkpChangeLogEntryV1(dcle))
}

func TxUpsertDKPChange(tx *bbolt.Tx, dcle record.DKPChangeEntry) error {
	return db.TxUpsert(tx, dcle.GetEntryId(), newDkpChangeLogEntryV1(dcle))
}

func TxGetDKPChangesForTarget(tx *bbolt.Tx, target string) ([]record.DKPChangeEntry, error) {
	records:=new([]record.DKPChangeEntry)
	err := db.TxForEach(tx, bolthold.Where("Target").Eq(target), func(entry *dkpChangeLogEntryV1)error {
		*records=append(*records, entry)
		return nil
	})
	sort.Sort(deltaByTimestamp(*records))
	return *records, err
}

type deltaByTimestamp []record.DKPChangeEntry
func (d deltaByTimestamp) Len() int {return len(d)}
func (d deltaByTimestamp) Less(i, j int) bool {return d[i].GetTimestamp().After(d[j].GetTimestamp())}
func (d deltaByTimestamp) Swap(i, j int) {d[i],d[j] = d[j],d[i]}

func TxRemoveDKPChange(tx *bbolt.Tx, entryId uint64) error {
	return db.TxDelete(tx, entryId, &dkpChangeLogEntryV1{})
}

func TxGetDKPChange(tx *bbolt.Tx, entryId uint64) (record.DKPChangeEntry, error) {
	var res dkpChangeLogEntryV1
	err := db.TxGet(tx, entryId, &res)
	return &res, err
}

func GetDKPChangesForTarget(target string) ([]record.DKPChangeEntry, error) {
	records:=new([]record.DKPChangeEntry)
	err := db.MakeView(func(tx *bbolt.Tx)error {
		var err error
		*records, err = TxGetDKPChangesForTarget(tx, target)
		return err
	})
	return *records, err
}

func GetDKPChanges() ([]record.DKPChangeEntry, error) {
	records:=new([]record.DKPChangeEntry)
	err := db.ForEach(&bolthold.Query{}, func (entry *dkpChangeLogEntryV1)error {
		*records=append(*records, entry)
		return nil
	})
	sort.Sort(deltaByTimestamp(*records))
	return *records, err
}

type dkpChangeLogEntryV1 struct {
	EntryId uint64 `boltholdKey:"EntryId"`
	Timestamp time.Time
	Target string
	Delta float64
	Description string

	RaidId    uint64
	Authority string
}

func (d *dkpChangeLogEntryV1) GetEntryId() uint64 {return d.EntryId}
func (d *dkpChangeLogEntryV1) GetTimestamp() time.Time {return d.Timestamp}
func (d *dkpChangeLogEntryV1) GetTarget() string {return d.Target}
func (d *dkpChangeLogEntryV1) GetDelta() float64 {return d.Delta}
func (d *dkpChangeLogEntryV1) GetDescription() string {return d.Description}
func (d *dkpChangeLogEntryV1) GetAuthority() string {return d.Authority}
func (d *dkpChangeLogEntryV1) GetRaidId() uint64 {return d.RaidId}

func newDkpChangeLogEntryV1(dce record.DKPChangeEntry) *dkpChangeLogEntryV1 {
	return &dkpChangeLogEntryV1{
		Timestamp:   dce.GetTimestamp(),
		Target:      dce.GetTarget(),
		Delta:       dce.GetDelta(),
		Description: dce.GetDescription(),
		RaidId:      dce.GetRaidId(),
		Authority:   dce.GetAuthority(),
	}
}
