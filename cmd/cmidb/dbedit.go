// +build server

package main

import (
	"encoding/json"
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"github.com/GontikR99/chillmodeinfo/internal/record"
	"os"
	"time"
)

func main() {
	if len(os.Args)<2 {
		fmt.Println("What do you want me to do?  list/promote/demote/listmembers/wipemembers/showlogs/showalllogs")
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