// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/timshannon/bolthold"
	"log"
)

type userProfileV0 struct {
	IdToken     string `boltholdKey:"IdToken"`
	Email       string
	DisplayName string
}

func RegisterUser(idToken string, email string) error {
	var profile userProfileV0
	err := db.Get(idToken, &profile)
	if err == nil {
		return nil
	}
	if err != bolthold.ErrNotFound {
		return err
	}
	profile = userProfileV0{
		IdToken:  idToken,
		Email:    email,
		DisplayName: "",
	}
	log.Println("Registering ", profile)
	return db.Upsert(idToken, &profile)
}