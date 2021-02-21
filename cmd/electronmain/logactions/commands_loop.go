// +build wasm,electron

package logactions

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"github.com/GontikR99/chillmodeinfo/internal/settings"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"regexp"
	"strings"
)

func init() {
	settings.DefaultSetting(settings.BidStartPattern, "^{C} tells you, '!start")
	settings.DefaultSetting(settings.BidEndPattern, "^{C} tells you, '!end")
	settings.DefaultSetting(settings.BidClosePattern, "^{C} tells you, '!clear")
}

var youSayRe=regexp.MustCompile("^You say, '(.*)'$")
var youOocRe=regexp.MustCompile("^You say out of character, '(.*)'$")
var youShoutRe=regexp.MustCompile("^You shout, '(.*)'$")
var youAuctionRe=regexp.MustCompile("^You auction, '(.*)'$")
var youGuildRe=regexp.MustCompile("^You say to your guild, '(.*)'$")
var youRaidRe=regexp.MustCompile("^You tell your raid, '(.*)'$")
var youGroupRe=regexp.MustCompile("^You tell your party, '(.*)'$")
var youChannelRe=regexp.MustCompile("^You tell [^':]*:[0-9]*, '(.*)'$")

var tellRe=regexp.MustCompile("^([^']*) tells you, '(.*)'$")

func selfMessage(entry *eqspec.LogEntry) string {
	if m:=youSayRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youShoutRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youAuctionRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youGuildRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=tellRe.FindStringSubmatch(entry.Message); m!=nil && strings.EqualFold(m[1], entry.Character) {
		return m[2]
	} else if m:=youOocRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youChannelRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youRaidRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else if m:=youGroupRe.FindStringSubmatch(entry.Message); m!=nil {
		return m[1]
	} else {
		return ""
	}
}

func matchesSetting(entry *eqspec.LogEntry, settingKey string) bool {
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
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		for i:=0;i<len(entries);i++ {
			entry := entries[i]

			selfMsg := selfMessage(entry)
			for _, item := range eqspec.BuiltTrie.Scan(selfMsg) {
				console.Log("Detected item: "+item)
			}

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