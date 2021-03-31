// +build server

package serverrpcs

import (
	"context"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/internal/util"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"log"
	"time"
	"unicode"
)

type serverRaidStub struct{}

func isValidCharacter(charname string) bool {
	for _, c := range charname {
		if !unicode.Is(unicode.Latin, c) {
			return false
		}
	}
	return charname!=""
}

func (s serverRaidStub) Fetch(ctx context.Context) ([]record.Raid, error) {
	raidsHolder := new([]record.Raid)
	err := db.MakeView([]db.TableName{dao.TableRaid, dao.TableDKPLog}, func(tx *bbolt.Tx) error {
		creditedByRaid := make(map[uint64]map[string]struct{})
		logs, err := dao.TxGetDKPChanges(tx)
		if err != nil {
			return nil
		}
		for _, log := range logs {
			if log.GetRaidId() == 0 {
				continue
			}
			if _, ok := creditedByRaid[log.GetRaidId()]; !ok {
				creditedByRaid[log.GetRaidId()] = make(map[string]struct{})
			}
			creditedByRaid[log.GetRaidId()][log.GetTarget()] = struct{}{}
		}

		raids, err := dao.TxGetRaids(tx)
		if err != nil {
			return err
		}
		for _, v := range raids {
			br := record.NewBasicRaid(v)
			br.Credited = nil
			if creditMap, ok := creditedByRaid[v.GetRaidId()]; ok {
				for k, _ := range creditMap {
					br.Credited = append(br.Credited, k)
				}
			}
			*raidsHolder = append(*raidsHolder, br)
		}
		return nil
	})
	return *raidsHolder, err
}

func (s serverRaidStub) Add(ctx context.Context, raid record.Raid, postDateHours int) error {
	log.Println("Starting add")
	selfProfile, err := requiresAdmin(ctx)
	if err != nil {
		return err
	}

	return db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog, dao.TableRaid}, func(tx *bbolt.Tx) error {
		newRaid := record.NewBasicRaid(raid)
		newRaid.RaidId = 0
		nowTime := time.Now()
		newRaid.Timestamp = nowTime.Add(time.Duration(postDateHours)*time.Hour)

		raidId, err := dao.TxInsertRaid(tx, newRaid)
		if err != nil {
			return err
		}

		seenMembers := make(map[string]struct{})
		for _, attendeeTmp := range raid.GetAttendees() {
			attendee := initialCap(attendeeTmp)
			if !isValidCharacter(attendee) {continue}
			if _, present := seenMembers[attendee]; present {
				continue
			}
			seenMembers[attendee] = struct{}{}
			member, err := dao.TxGetMember(tx, attendee)
			if err == bolthold.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}
			if member.IsAlt() {
				continue
			}

			newMember := record.NewBasicMember(member)
			if newMember.LastActive.Before(nowTime) {
				newMember.LastActive = nowTime
			}
			newMember.DKP += raid.GetDKPValue()
			err = dao.TxUpsertMember(tx, newMember)
			if err != nil {
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
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s serverRaidStub) Delete(ctx context.Context, raidId uint64) error {
	_, err := requiresAdmin(ctx)
	if err != nil {
		return err
	}

	return db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog, dao.TableRaid}, func(tx *bbolt.Tx) error {
		_, err := dao.TxGetRaid(tx, raidId)
		if err != nil {
			return err
		}
		err = dao.TxDeleteRaid(tx, raidId)
		if err != nil {
			return err
		}

		changeSet, err := dao.TxGetDKPChangesForRaid(tx, raidId)
		if err != nil {
			return err
		}

		for _, changeEntry := range changeSet {
			err = dao.TxRemoveDKPChange(tx, changeEntry.GetEntryId())
			if err != nil {
				return err
			}

			member, err := dao.TxGetMember(tx, changeEntry.GetTarget())
			if err == bolthold.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}
			newMember := record.NewBasicMember(member)
			newMember.DKP -= changeEntry.GetDelta()
			err = dao.TxUpsertMember(tx, newMember)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s serverRaidStub) Update(ctx context.Context, raid record.Raid) (record.Raid, error) {
	selfProfile, err := requiresAdmin(ctx)
	if err != nil {
		return nil, err
	}

	newRaidHolder := new(record.Raid)
	err = db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog, dao.TableRaid}, func(tx *bbolt.Tx) error {
		activeTime := raid.GetTimestamp()
		nowTime := time.Now()
		if activeTime.After(nowTime) {
			activeTime = nowTime
		}
		oldRaid, err := dao.TxGetRaid(tx, raid.GetRaidId())
		if err != nil {
			return err
		}

		newRaid := record.NewBasicRaid(oldRaid)
		newRaid.Description = raid.GetDescription()
		newRaid.DKPValue = raid.GetDKPValue()
		newRaid.Attendees = nil
		for _, a := range raid.GetAttendees() {
			if isValidCharacter(a) {
				newRaid.Attendees = append(newRaid.Attendees, initialCap(a))
			}
		}
		newRaid.Attendees = util.Deduplicate(newRaid.Attendees)
		newRaid.Credited = nil

		err = dao.TxUpsertRaid(tx, newRaid)
		if err != nil {
			return err
		}
		*newRaidHolder = newRaid

		// Remove old changes
		changeSet, err := dao.TxGetDKPChangesForRaid(tx, raid.GetRaidId())
		if err != nil {
			return err
		}

		for _, changeEntry := range changeSet {
			err = dao.TxRemoveDKPChange(tx, changeEntry.GetEntryId())
			if err != nil {
				return err
			}

			member, err := dao.TxGetMember(tx, changeEntry.GetTarget())
			if err == bolthold.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}
			newMember := record.NewBasicMember(member)
			newMember.DKP -= changeEntry.GetDelta()
			err = dao.TxUpsertMember(tx, newMember)
			if err != nil {
				return err
			}
		}

		// create new changes
		seenMembers := make(map[string]struct{})
		for _, attendeeTmp := range util.Deduplicate(raid.GetAttendees()) {
			attendee := initialCap(attendeeTmp)
			if !isValidCharacter(attendee) {continue}
			if _, present := seenMembers[attendee]; present {
				continue
			}
			seenMembers[attendee] = struct{}{}
			member, err := dao.TxGetMember(tx, attendee)
			if err == bolthold.ErrNotFound {
				continue
			} else if err != nil {
				return err
			}
			if member.IsAlt() {
				continue
			}

			newMember := record.NewBasicMember(member)
			if newMember.LastActive.Before(activeTime) {
				newMember.LastActive = activeTime
			}
			newMember.DKP += raid.GetDKPValue()
			err = dao.TxUpsertMember(tx, newMember)
			if err != nil {
				return err
			}

			newLogEntry := &record.BasicDKPChangeEntry{
				EntryId:     0,
				Timestamp:   newRaid.GetTimestamp(),
				Target:      attendee,
				Delta:       newRaid.GetDKPValue(),
				Description: newRaid.Description,
				RaidId:      newRaid.GetRaidId(),
				Authority:   selfProfile.GetDisplayName(),
			}

			err = dao.TxAppendDKPChange(tx, newLogEntry)
			if err != nil {
				return err
			}

			newRaid.Credited = append(newRaid.Credited, attendee)
		}
		return nil
	})
	return *newRaidHolder, err
}

func init() {
	register(restidl.HandleRaid(&serverRaidStub{}))
}
