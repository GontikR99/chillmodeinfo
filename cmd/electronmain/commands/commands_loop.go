// +build wasm,electron

package commands

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"regexp"
	"strings"
)

func init() {
	settings.DefaultSetting(settings.BidStartPattern, "^You say to your guild, 'BIDS START")
	settings.DefaultSetting(settings.BidEndPattern, "^You say to your guild, 'BIDS END")
	settings.DefaultSetting(settings.BidClosePattern, "^{C} tells you, '!clear")
}

func matchesSetting(entry *eqfiles.LogEntry, settingKey string) bool {
	pattern, present, err := settings.LookupSetting(settingKey)
	if err!=nil || !present {
		return false
	}
	pattern = strings.ToUpper(pattern)
	pattern = strings.ReplaceAll(pattern, "{C}", strings.ToUpper(entry.Character))
	re, err := regexp.Compile(pattern)
	if err!=nil {
		return false
	}
	return re.MatchString(strings.ToUpper(entry.Message))
}

func WatchLogs() {
	eqfiles.RegisterLogsListener(func(entries []*eqfiles.LogEntry) {
		for i:=0;i<len(entries);i++ {
			entry := entries[i]
			if matchesSetting(entry, settings.BidStartPattern) {
				OpenBids(entries[i+1:])
			}
			if matchesSetting(entry, settings.BidEndPattern) {
				EndBids()
			}
			if matchesSetting(entry, settings.BidClosePattern) {
				ClearBids()
			}
		}
	})
}