// +build wasm

package localprofile

import "github.com/GontikR99/chillmodeinfo/internal/profile"

var currentProfile profile.Entry

func SetProfileIfAbsent(entry profile.Entry) {
	if currentProfile==nil && entry!=nil {
		currentProfile=entry
	}
}

func GetProfile() profile.Entry {
	return currentProfile
}

const channelProfile="profiles"

type profileMessage struct {
	Value *profile.BasicProfile
}

func IsAdmin() bool {
	return GetProfile()!=nil && GetProfile().GetAdminState()==profile.StateAdminApproved
}