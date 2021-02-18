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
	return db.TxInsert(tx, bolthold.NextSequence(), newDkpChangeLogEntryV0(dcle))
}

func TxGetDKPChangesForTarget(tx *bbolt.Tx, target string) ([]record.DKPChangeEntry, error) {
	records:=new([]record.DKPChangeEntry)
	err := db.TxForEach(tx, bolthold.Where("Target").Eq(target), func(entry *dkpChangeLogEntryV0)error {
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
	err := db.ForEach(bolthold.Where("Target").Ne(""), func (entry *dkpChangeLogEntryV0)error {
		*records=append(*records, entry)
		return nil
	})
	sort.Sort(deltaByTimestamp(*records))
	return *records, err
}

type dkpChangeLogEntryV0 struct {
	Timestamp time.Time
	Target string
	Delta float64
	Description string

	RaidId    uint64
	Authority string
}

func (dcle *dkpChangeLogEntryV0) GetTimestamp() time.Time {return dcle.Timestamp}
func (dcle *dkpChangeLogEntryV0) GetTarget() string {return dcle.Target}
func (dcle *dkpChangeLogEntryV0) GetDelta() float64 {return dcle.Delta}
func (dcle *dkpChangeLogEntryV0) GetDescription() string {return dcle.Description}
func (dcle *dkpChangeLogEntryV0) GetAuthority() string {return dcle.Authority}
func (dcle *dkpChangeLogEntryV0) GetRaidId() uint64    {return dcle.RaidId }

func newDkpChangeLogEntryV0(dce record.DKPChangeEntry) *dkpChangeLogEntryV0 {
	return &dkpChangeLogEntryV0{
		Timestamp:   dce.GetTimestamp(),
		Target:      dce.GetTarget(),
		Delta:       dce.GetDelta(),
		Description: dce.GetDescription(),
		RaidId:      dce.GetRaidId(),
		Authority:   dce.GetAuthority(),
	}
}