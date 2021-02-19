// +build server

package serverrpcs

import (
	"context"
	"errors"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"go.etcd.io/bbolt"
	"log"
	"time"
)

type serverRaidStub struct {}

func (s serverRaidStub) Fetch(ctx context.Context) ([]record.Raid, error) {
	raidsHolder:=new([]record.Raid)
	err := db.MakeView(func(tx *bbolt.Tx) error {
		raids, err := dao.TxGetRaids(tx)
		if err!=nil {return err}
		for _, v := range raids {
			*raidsHolder = append(*raidsHolder, record.NewBasicRaid(v))
		}
		return nil
	})
	return *raidsHolder, err
}


func (s serverRaidStub) Add(ctx context.Context, raid record.Raid) error {
	log.Println("Starting add")
	selfProfile, err := requiresAdmin(ctx)
	if err!=nil {return err}

	return db.MakeUpdate(func(tx *bbolt.Tx) error {
		newRaid := record.NewBasicRaid(raid)
		newRaid.RaidId = 0
		newRaid.Timestamp = time.Now()

		raidId, err := dao.TxInsertRaid(tx, newRaid)
		if err!=nil {
			return err
		}

		seenMembers := make(map[string]struct{})
		for _, attendeeTmp := range raid.GetAttendees() {
			attendee := initialCap(attendeeTmp)
			if _, present := seenMembers[attendee]; present {
				continue
			}
			seenMembers[attendee]=struct{}{}
			member, err := dao.TxGetMember(tx, attendee)
			if err!=nil {
				return err
			}
			if member.IsAlt() {
				continue
			}

			newMember:=record.NewBasicMember(member)
			newMember.LastActive=time.Now()
			newMember.DKP += raid.GetDKPValue()
			err = dao.TxUpsertMember(tx, newMember)
			if err!=nil {
				return err
			}

			newLogEntry := &record.BasicDKPChangeEntry{
				EntryId:     0,
				Timestamp:   newRaid.GetTimestamp(),
				Target:      attendee,
				Delta:       newRaid.GetDKPValue(),
				Description: newRaid.Description,
				RaidId:      raidId,
				Authority:   selfProfile.GetDisplayName(),
			}

			err = dao.TxAppendDKPChange(tx, newLogEntry)
			if err!=nil {
				return err
			}
		}
		return nil
	})
}

func (s serverRaidStub) Delete(ctx context.Context, raidId uint64) error {
	_, err := requiresAdmin(ctx)
	if err!=nil {return err}

	return db.MakeUpdate(func(tx *bbolt.Tx) error {
		_, err := dao.TxGetRaid(tx, raidId)
		if err!=nil {
			return err
		}
		err = dao.TxDeleteRaid(tx, raidId)
		if err!=nil {
			return err
		}

		changeSet, err := dao.TxGetDKPChangesForRaid(tx, raidId)
		if err!=nil {
			return err
		}

		for _, changeEntry := range changeSet {
			err = dao.TxRemoveDKPChange(tx, changeEntry.GetEntryId())
			if err!=nil {return err}

			member, err := dao.TxGetMember(tx, changeEntry.GetTarget())
			if err!=nil {
				return err
			}
			newMember := record.NewBasicMember(member)
			newMember.DKP -= changeEntry.GetDelta()
			err = dao.TxUpsertMember(tx, newMember)
			if err!=nil {
				return err
			}
		}
		return nil
	})
}

func (s serverRaidStub) Update(ctx context.Context, raid record.Raid) (record.Raid, error) {
	return nil, errors.New("unsupported")
}

func init() {
	register(restidl.HandleRaid(&serverRaidStub{}))
}