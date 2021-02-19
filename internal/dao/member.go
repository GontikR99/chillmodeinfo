// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"time"
)

const TableMembers=db.TableName("Members")

type memberV0 struct {
	Name string `boltholdKey:"Name"`
	Class string
	Level int16
	Rank string
	Alt bool
	DKP float64
	LastActive time.Time
	Owner string
}

func newMemberV0(m record.Member) *memberV0 {
	return &memberV0{
		Name:       m.GetName(),
		Class:      m.GetClass(),
		Level:      m.GetLevel(),
		Rank: m.GetRank(),
		Alt:        m.IsAlt(),
		DKP:        m.GetDKP(),
		LastActive: m.GetLastActive(),
		Owner: m.GetOwner(),
	}
}

func (m *memberV0) GetName() string {return m.Name}
func (m *memberV0) GetClass() string {return m.Class}
func (m *memberV0) GetLevel() int16 {return m.Level}
func (m *memberV0) GetRank() string {return m.Rank}
func (m *memberV0) IsAlt() bool {return m.Alt}
func (m *memberV0) GetDKP() float64 {return m.DKP}
func (m *memberV0) GetLastActive() time.Time {return m.LastActive}
func (m *memberV0) GetOwner() string {return m.Owner}

func TxGetMembers(tx *bbolt.Tx) (map[string]record.Member, error) {
	result:=make(map[string]record.Member)
	db.TxForEach(tx, bolthold.Where("Name").Ne(""), func(m *memberV0) error {
		result[m.GetName()]=record.NewBasicMember(m)
		return nil
	})
	return result, nil
}

func GetMembers() (map[string]record.Member, error) {
	result:=new(map[string]record.Member)
	err := db.MakeView([]db.TableName{TableMembers}, func(tx *bbolt.Tx) error {
		var err error
		*result, err = TxGetMembers(tx)
		return err
	})
	return *result, err
}

func TxGetMember(tx *bbolt.Tx, name string) (record.Member, error) {
	var m memberV0
	err := db.TxGet(tx, name, &m)
	return &m, err
}

func GetMember(name string) (record.Member, error) {
	m:=new(record.Member)
	err := db.MakeView([]db.TableName{TableMembers}, func(tx *bbolt.Tx) error {
		var err error
		*m, err = TxGetMember(tx, name)
		return err
	})
	return *m, err
}

func TxUpsertMember(tx *bbolt.Tx, m record.Member) error {
	return db.TxUpsert(tx, m.GetName(), newMemberV0(m))
}

func UpsertMember(m record.Member) error {
	return db.MakeUpdate([]db.TableName{TableMembers}, func(tx *bbolt.Tx) error {
		return TxUpsertMember(tx, m)
	})
}

func WipeMembers() error {
	return db.MakeUpdate([]db.TableName{TableMembers}, func(tx *bbolt.Tx) error {
		return db.TxDeleteMatching(tx, &memberV0{}, bolthold.Where("Name").Ne(""))
	})
}