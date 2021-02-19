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
	"os"
	"time"
)

func main() {
	if len(os.Args)<2 {
		fmt.Println("What do you want me to do?  list/promote/demote/listmembers/wipemembers/showlogs/showalllogs/zeromember")
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
	case "wipemembers":
		fmt.Println("Wiping members.")
		dao.WipeMembers()
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

	default:
		fmt.Println("Unknown command ", os.Args[1])
	}
}