// +build server

package main

import (
	"encoding/json"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/dao/db"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const syncText="Sync from Gamerlaunch"
const syncAuth="Gamerlaunch"
const syncPage="https://chillmode.gamerlaunch.com/rapid_raid/leaderboard.php"

func main() {
	if len(os.Args)<2 {
		fmt.Println("What do you want me to do?  list/promote/demote/listmembers/wipemembers/showlogs/showalllogs/syncgamerlaunch/zeromember")
	}
	switch os.Args[1] {
	case "list":
		for _, v :=range dao.ListAllProfiles() {
			fmt.Println("UserId: ", v.GetUserId())
			fmt.Println("DisplayName: ", v.GetDisplayName())
			fmt.Println("Email: ", v.GetEmail())
			fmt.Println("AdminState: ", v.GetAdminState())
			fmt.Println("StartTime: ", v.GetStartDate().Format(time.ANSIC))
			fmt.Println()
		}
	case "promote":
		if len(os.Args)<4 {
			fmt.Println("promote <id> <display-name>")
			return
		}
		err := dao.UpdateProfileForAdmin(os.Args[2], os.Args[3], profile.StateAdminApproved, time.Unix(0,0))
		if err!=nil {
			fmt.Println(err)
		} else {
			fmt.Println("Updated %s(%s)", os.Args[3], os.Args[2])
		}
	case "demote":
		if len(os.Args)<3 {
			fmt.Println("demote <id>")
			return
		}
		err := dao.UpdateProfileForAdmin(os.Args[2], "", profile.StateAdminUnrequested, time.Unix(0,0))
		if err!=nil {
			fmt.Println(err)
		} else {
			fmt.Println("Updated %s", os.Args[2])
		}
	case "listmembers":
		fmt.Println("Listing members.")
		m, _ := dao.GetMembers()
		for k, v:=range m {
			fmt.Printf("%v = %v\n", k, *record.NewBasicMember(v))
		}
	case "wipe":
		fmt.Println("Wiping members/etc.")
		dao.WipeMembers()
		dao.WipeDKP()
		dao.WipeRaids()

	case "zeromember":
		if len(os.Args)<3 {
			fmt.Println("zeromember <name>")
		}
		name:=os.Args[2]
		err:=db.MakeUpdate([]db.TableName{dao.TableMembers, dao.TableDKPLog},func(tx *bbolt.Tx) error {
			member, err := dao.GetMember(name)
			if err!=nil {
				return err
			}
			newMember := record.NewBasicMember(member)
			newMember.DKP=0
			newMember.LastActive=time.Time{}

			dkpEntries, err := dao.GetDKPChangesForTarget(name)
			if err!=nil {
				return err
			}
			for _, entry := range dkpEntries {
				err = dao.TxRemoveDKPChange(tx, entry.GetEntryId())
				if err!=nil {return err}
			}
			return dao.TxUpsertMember(tx, newMember)
		})
		if err!=nil {
			panic(err)
		}
	case "showlogs":
		if len(os.Args)<3 {
			fmt.Println("showlogs <id>")
			return
		}
		fmt.Println("Listing logs for "+os.Args[2])
		entries, err := dao.GetDKPChangesForTarget(os.Args[2])
		if err!=nil {
			panic(err)
		}
		for idx, entry:=range entries {
			jj, _ := json.Marshal(record.NewBasicDKPChangeEntry(entry))
			fmt.Printf("%d=%s\n", 1+idx, string(jj))
		}
		fmt.Println("done")

	case "showalllogs":
		fmt.Println("Listing all DKP logs")
		entries, err := dao.GetDKPChanges()
		if err!=nil {
			panic(err)
		}
		for idx, entry:=range entries {
			jj, _ := json.Marshal(record.NewBasicDKPChangeEntry(entry))
			fmt.Printf("%d=%s\n", 1+idx, string(jj))
		}
		fmt.Println("done")

	case "syncgamerlaunch":
		fmt.Println("Syncing from gamerlaunch")
		appends := 0
		updates := 0
		resp, err := http.Get(syncPage)
		if err!=nil {
			panic(err)
		}
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err!=nil {
			panic(err)
		}
		bodyText := string(bodyBytes)
		entryRe:=regexp.MustCompile("<tr><td[^>]*><a href='/users[^>]*>([A-Za-z]*)</a></td><td[^>]*><span class='dkp_current'>([0-9.,]*)</span></td></tr>")
		matches := entryRe.FindAllStringSubmatch(bodyText, -1)
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
			panic(err)
		}
		fmt.Printf("%d appends, %d updates\n", appends, updates)
	default:
		fmt.Println("Unknown command ", os.Args[1])
	}
}