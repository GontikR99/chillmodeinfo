// +build server

package main

import (
	"fmt"
	"github.com/GontikR99/chillmodeinfo/internal/dao"
	"github.com/GontikR99/chillmodeinfo/internal/profile"
	"os"
	"time"
)

func main() {
	if len(os.Args)<2 {
		fmt.Println("What do you want me to do?  list/promote/demote")
	}
	switch os.Args[1] {
	case "list":
		for _, v :=range dao.ListAllProfiles() {
			fmt.Println("UserId: ", v.GetUserId())
			fmt.Println("DisplayName: ", v.GetDisplayName())
			fmt.Println("Email: ", v.GetEmail())
			fmt.Println("AdminState: ", v.GetAdminState())
			fmt.Println("StartTime: ", v.GetAdminStartDate().Format(time.ANSIC))
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
	default:
		fmt.Println("Unknown command ", os.Args[1])
	}
}