// +build server

package dao

import (
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/timshannon/bolthold"
)

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
		UserId: idToken,
		Email:  email,
	}
	return db.Upsert(idToken, &profile)
}

type userProfileV0 struct {
	UserId      string `boltholdKey:"UserId"`
	Email       string
	DisplayName string
	AdminState  profile.AdminState
}

func (u *userProfileV0) GetIdToken() string  {return u.UserId }
func (u *userProfileV0) GetEmail() string               {return u.Email }
func (u *userProfileV0) GetDisplayName() string            {return u.DisplayName }
func (u *userProfileV0) GetAdminState() profile.AdminState {return u.AdminState }

func LookupProfile(userId string) (profile.Entry, error) {
	var entry userProfileV0
	err := db.Get(userId, &entry)
	return &entry, err
}