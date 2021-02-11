// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/timshannon/bolthold"
)

type clientIdAssociationV0 struct {
	ClientId string `boltholdKey:"ClientId"`
	GoogleId string
}

func LookupClientId(clientId string) (string, bool, error) {
	ciaV0 := &clientIdAssociationV0{}
	err := db.Get(clientId, ciaV0)
	if err == bolthold.ErrNotFound {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	} else {
		return ciaV0.GoogleId, true, nil
	}
}

func AssociateClientId(clientId string, googleId string) error {
	ciaV0 := &clientIdAssociationV0{
		ClientId: clientId,
		GoogleId: googleId,
	}
	return db.Upsert(clientId, ciaV0)
}

func DisassociateClientId(clientId string) error {
	return db.Delete(clientId, new(clientIdAssociationV0))
}
