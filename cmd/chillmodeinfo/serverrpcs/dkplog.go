// +build server

package serverrpcs

import (
	"context"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/comms/httputil"
	"github.com/GontikR99/chillmodeinfo/internal/comms/restidl"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"github.com/GontikR99/chillmodeinfo/internal/sitedef"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type serverDKPLogHandler struct{}

func (s serverDKPLogHandler) Append(ctx context.Context, delta record.DKPChangeEntry) error {
	selfProfile, err := requiresAdmin(ctx)
	if err != nil {
		return err
	}
	return db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog}, func(tx *bbolt.Tx) error {
		target := initialCap(delta.GetTarget())
		targetMemberRecord, err := dao.TxGetMember(tx, target)
		if err != nil {
			return err
		}
		nowTime := time.Now()
		updatedMember := record.NewBasicMember(targetMemberRecord)
		updatedMember.LastActive = nowTime
		updatedMember.DKP = updatedMember.DKP + delta.GetDelta()

		updatedDelta := record.NewBasicDKPChangeEntry(delta)
		updatedDelta.Target = updatedMember.GetName()
		updatedDelta.Authority = selfProfile.GetDisplayName()
		updatedDelta.Timestamp = nowTime

		err = dao.TxUpsertMember(tx, updatedMember)
		if err != nil {
			return err
		}
		return dao.TxAppendDKPChange(tx, updatedDelta)
	})
}

func (s serverDKPLogHandler) Retrieve(ctx context.Context, target string) ([]record.DKPChangeEntry, error) {
	if target != "" {
		return dao.GetDKPChangesForTarget(target)
	} else {
		all, err := dao.GetDKPChanges()
		if err != nil {
			return nil, err
		}
		var filtered []record.DKPChangeEntry
		for _, v := range all {
			if v.GetRaidId() == 0 {
				filtered = append(filtered, v)
			}
		}
		return filtered, nil
	}
}

func (s serverDKPLogHandler) Remove(ctx context.Context, entryId uint64) error {
	_, err := requiresAdmin(ctx)
	if err != nil {
		return err
	}

	return db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog}, func(tx *bbolt.Tx) error {
		dkpEntry, err := dao.TxGetDKPChange(tx, entryId)
		if err != nil {
			return err
		}
		if dkpEntry.GetRaidId() != 0 {
			return httputil.NewError(http.StatusBadRequest, "You may not remove a log entry corresponding to a raid.  Go remove the raid.")
		}
		dao.TxRemoveDKPChange(tx, entryId)

		member, err := dao.TxGetMember(tx, dkpEntry.GetTarget())
		if err != nil {
			return err
		}

		newMember := record.NewBasicMember(member)
		newMember.DKP -= dkpEntry.GetDelta()
		newMember.LastActive = time.Now()

		return dao.TxUpsertMember(tx, newMember)
	})
}

func (s serverDKPLogHandler) Update(ctx context.Context, newEntry record.DKPChangeEntry) (record.DKPChangeEntry, error) {
	selfProfile, err := requiresAdmin(ctx)
	if err != nil {
		return nil, err
	}

	resHolder := new(record.DKPChangeEntry)
	err = db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog}, func(tx *bbolt.Tx) error {
		oldEntry, err := dao.TxGetDKPChange(tx, newEntry.GetEntryId())
		if err != nil {
			return err
		}
		if oldEntry.GetRaidId() != 0 {
			return httputil.NewError(http.StatusBadRequest, "You may not update a log entry corresponding to a raid.  Go update the raid.")
		}
		if oldEntry.GetTarget() != newEntry.GetTarget() {
			return httputil.NewError(http.StatusBadRequest, "You may not change the target of a DKP update with this API.")
		}
		newEntry := record.NewBasicDKPChangeEntry(newEntry)
		*resHolder = newEntry
		newEntry.Timestamp = oldEntry.GetTimestamp()
		newEntry.Authority = selfProfile.GetDisplayName()
		err = dao.TxUpsertDKPChange(tx, newEntry)
		if err != nil {
			return err
		}

		member, err := dao.TxGetMember(tx, newEntry.GetTarget())
		if err != nil {
			return err
		}

		newMember := record.NewBasicMember(member)
		newMember.DKP += newEntry.GetDelta() - oldEntry.GetDelta()

		return dao.TxUpsertMember(tx, newMember)
	})
	return *resHolder, err
}

const syncText="Sync from Gamerlaunch"
const syncAuth="Gamerlaunch"
var gamerlaunchScrapeRe=regexp.MustCompile("<tr><td[^>]*><a href='/users[^>]*>([A-Za-z]*)</a></td><td[^>]*><span class='dkp_current'>([0-9.,]*)</span></td></tr>")

func (s serverDKPLogHandler) Sync(ctx context.Context) (string, error) {
	_, err := requiresAdmin(ctx)
	if err!=nil {
		return "", err
	}

	appends := 0
	updates := 0
	resp, err := http.Get(sitedef.GamerlaunchSyncURL)
	if err!=nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		panic(err)
	}
	bodyText := string(bodyBytes)
	matches := gamerlaunchScrapeRe.FindAllStringSubmatch(bodyText, -1)
	syncDKP := make(map[string]float64)
	for _, m := range matches {
		dkpString := strings.ReplaceAll(m[2], ",", "")
		dkpValue, err := strconv.ParseFloat(dkpString, 64)
		if err!=nil {
			panic(err)
		}
		syncDKP[m[1]]=dkpValue
	}

	err = db.MakeUpdate([]db.TableName{dao.TableDKPLog, dao.TableMembers}, func(tx *bbolt.Tx) error {
		members, err := dao.TxGetMembers(tx)
		if err!=nil {
			return err
		}
		logs, err := dao.TxGetDKPChanges(tx)
		if err!=nil {
			return err
		}
		// Update old sync entries
		for _, logEntry := range logs {
			if _, present := syncDKP[logEntry.GetTarget()]; !present {
				continue
			}
			if _, present := members[logEntry.GetTarget()]; !present {
				continue
			}
			if logEntry.GetDescription()==syncText && logEntry.GetAuthority()==syncAuth {
				updateEntry := record.NewBasicDKPChangeEntry(logEntry)
				change := syncDKP[logEntry.GetTarget()] - logEntry.GetDelta()
				delete(syncDKP, logEntry.GetTarget())
				if change==0 {
					continue
				}
				updateEntry.Delta += change
				err = dao.TxUpsertDKPChange(tx, updateEntry)
				if err!=nil {
					return err
				}

				updateMember := record.NewBasicMember(members[logEntry.GetTarget()])
				if updateMember.LastActive.Before(updateEntry.Timestamp) {
					updateMember.LastActive = updateEntry.Timestamp
				}
				updateMember.DKP += change
				err = dao.TxUpsertMember(tx, updateMember)
				if err!=nil {
					return err
				}
				updates++
			}
		}

		// Insert new sync entries
		for memberName, gamerlaunchDKP := range syncDKP {
			if _, present := members[memberName]; !present {
				continue
			}
			newEntry := &record.BasicDKPChangeEntry{
				EntryId:     0,
				Timestamp:   time.Now(),
				Target:      memberName,
				Delta:       gamerlaunchDKP,
				Description: syncText,
				RaidId:      0,
				Authority:   syncAuth,
			}
			err = dao.TxAppendDKPChange(tx, newEntry)
			if err!=nil {
				return err
			}

			updateMember := record.NewBasicMember(members[memberName])
			if updateMember.LastActive.Before(newEntry.Timestamp) {
				updateMember.LastActive=newEntry.Timestamp
			}
			updateMember.DKP += gamerlaunchDKP
			err = dao.TxUpsertMember(tx, updateMember)
			if err!=nil {
				return err
			}
			appends++
		}
		return nil
	})
	if err!=nil {
		return "", err
	}
	return fmt.Sprintf("%d new sync entries created.  %d old sync entries updated.\n", appends, updates), nil
}

func init() {
	register(restidl.HandleDKPLog(&serverDKPLogHandler{}))
}
