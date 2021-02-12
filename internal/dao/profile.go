// +build server

package dao

import (
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/timshannon/bolthold"
	"log"
	"time"
)

func RegisterUser(idToken string, email string) error {
	var profileObj userProfileV1
	err := db.Get(idToken, &profileObj)
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
	return db.Upsert(idToken, &profileObj)
}

func LookupProfile(userId string) (profile.Entry, error) {
	var entry userProfileV1
	err := db.Get(userId, &entry)
	return &entry, err
}

func UpdateProfileForAdmin(userId string, displayName string, adminState profile.AdminState, startDate time.Time) error {
	oldEntry, err := LookupProfile(userId)
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
	return db.Upsert(userId, &profileObj)
}

func ListAllProfiles() []profile.Entry {
	entries:=[]userProfileV1{}
	db.Find(&entries, bolthold.Where("UserId").Ne(""))
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

func init() {
	db.ForEach(bolthold.Where("UserId").Ne(""), func(record *userProfileV0) {
		log.Printf("userProfileV0 -> userProfileV1: %s", record.Email)
		db.Upsert(record.UserId, &userProfileV1{
			UserId:         record.UserId,
			Email:          record.Email,
			DisplayName:    record.DisplayName,
			AdminState:     record.AdminState,
			AdminStartDate: time.Time{},
		})
	})
	db.DeleteMatching(&userProfileV0{}, bolthold.Where("UserId").Ne(""))
}

type userProfileV0 struct {
	UserId      string `boltholdKey:"UserId"`
	Email       string
	DisplayName string
	AdminState  profile.AdminState
}