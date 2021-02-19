// +build server

package dao

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"time"
)

const TableProfile=db.TableName("Profile")

func RegisterUser(idToken string, email string) error {
	return db.MakeUpdate([]db.TableName{TableProfile}, func(tx *bbolt.Tx) error {
		var profileObj userProfileV1
		err := db.TxGet(tx, idToken, &profileObj)
		if err == nil {
			return nil
		}
		if err != bolthold.ErrNotFound {
			return err
		}
		profileObj = userProfileV1{
			UserId:     idToken,
			Email:      email,
			AdminState: profile.StateAdminUnrequested,
		}
		return db.TxUpsert(tx, idToken, &profileObj)
	})
}

func LookupProfile(userId string) (profile.Entry, error) {
	var entry userProfileV1
	err := db.MakeView([]db.TableName{TableProfile}, func(tx *bbolt.Tx) error {
		return db.TxGet(tx, userId, &entry)
	})
	return &entry, err
}

func UpdateProfileForAdmin(userId string, displayName string, adminState profile.AdminState, startDate time.Time) error {
	return db.MakeUpdate([]db.TableName{TableProfile}, func(tx *bbolt.Tx) error {
		var oldEntry userProfileV1
		err := db.TxGet(tx, userId, &oldEntry)
		if err!=nil {
			return err
		}
		if adminState<0 || adminState>=profile.StateAdminMax {
			return errors.New("Unsupported admin state")
		}
		profileObj := userProfileV1{
			UserId:         oldEntry.GetUserId(),
			Email:          oldEntry.GetEmail(),
			DisplayName:    displayName,
			AdminState:     adminState,
			AdminStartDate: startDate,
		}
		return db.TxUpsert(tx, userId, &profileObj)
	})
}

func ListAllProfiles() []profile.Entry {
	entries:=[]userProfileV1{}
	db.MakeView([]db.TableName{TableProfile}, func(tx *bbolt.Tx) error {
		return db.TxFind(tx, &entries, bolthold.Where("UserId").Ne(""))
	})
	retEntries:=[]profile.Entry{}
	for i:=0;i<len(entries); i++ {
		retEntries = append(retEntries, &entries[i])
	}
	return retEntries
}

type userProfileV1 struct {
	UserId      string `boltholdKey:"UserId"`
	Email       string
	DisplayName string
	AdminState  profile.AdminState
	AdminStartDate	time.Time
}

func (u *userProfileV1) GetUserId() string {return u.UserId}
func (u *userProfileV1) GetEmail() string {return u.Email}
func (u *userProfileV1) GetDisplayName() string {return u.DisplayName}
func (u *userProfileV1) GetAdminState() profile.AdminState {return u.AdminState}
func (u *userProfileV1) GetStartDate() time.Time           {return u.AdminStartDate}
