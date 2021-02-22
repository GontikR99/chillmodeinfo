// +build wasm,electron

package exerpcs

import (
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/internal/eqspec"
	"regexp"
	"strings"
)

type serverBidSupport struct {}

func (s *serverBidSupport) GetLastMentioned() (string, error) {
	return lastMentionedItem, nil
}

var bidOfferedFunc func(string, string, float64)
func OnBidOffered(bof func(string, string, float64)) {
	bidOfferedFunc=bof
}

func (s *serverBidSupport) OfferBid(bidder string, itemname string, bidValue float64) error {
	if bidOfferedFunc!=nil {
		bidOfferedFunc(bidder, itemname, bidValue)
	}
	return nil
}

var lastMentionedItem string

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

func init() {
	eqspec.RegisterLogsListener(func(entries []*eqspec.LogEntry) {
		for _, entry := range entries {
			selfMsg := selfMessage(entry)
			items := eqspec.BuiltTrie.Scan(selfMsg)
			if items!=nil {
				lastMentionedItem=items[0]
			}
		}
	})

	register(rpcidl.HandleBidSupport(&serverBidSupport{}))
}