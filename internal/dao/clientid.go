// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

type clientIdAssociationV0 struct {
	ClientId string `boltholdKey:"ClientId"`
	GoogleId string
}

const TableClientId=db.TableName("ClientId")

func LookupClientId(clientId string) (string, bool, error) {
	valPtr := new(string)
	presPtr := new(bool)
	err := db.MakeView([]db.TableName{TableClientId}, func(tx *bbolt.Tx) error {
		ciaV0 := &clientIdAssociationV0{}
		err := db.TxGet(tx, clientId, ciaV0)
		if err == bolthold.ErrNotFound {
			*presPtr=false
			return nil
		} else if err != nil {
			*presPtr=false
			return err
		} else {
			*valPtr=ciaV0.GoogleId
			*presPtr=true
			return nil
		}
	})
	return *valPtr, *presPtr, err
}

func AssociateClientId(clientId string, googleId string) error {
	ciaV0 := &clientIdAssociationV0{
		ClientId: clientId,
		GoogleId: googleId,
	}
	return db.MakeUpdate([]db.TableName{TableClientId}, func(tx *bbolt.Tx) error {
		return db.TxUpsert(tx, clientId, ciaV0)
	})
}

func DisassociateClientId(clientId string) error {
	return db.MakeUpdate([]db.TableName{TableClientId}, func(tx *bbolt.Tx) error {
		return db.TxDelete(tx, clientId, new(clientIdAssociationV0))
	})
}
