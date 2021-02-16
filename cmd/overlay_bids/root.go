// +build wasm,web

package main

import (
	"github.com/GontikR99/chillmodeinfo/internal/eqfiles"
	"github.com/GontikR99/chillmodeinfo/internal/comms/rpcidl"
	"github.com/GontikR99/chillmodeinfo/pkg/console"
	"github.com/GontikR99/chillmodeinfo/pkg/electron/ipc/ipcrenderer"
	"github.com/vugu/vugu"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Root struct {
	ActiveBids []*Bid
	RandBag    map[int]struct{}
}

type Bid struct {
	Name       string
	Value      int
	Tiebreaker int
	Texts      []string
}

func (c *Root) Init(ctx vugu.InitCtx) {
	rand.Seed(time.Now().UnixNano())
	c.RandBag = make(map[int]struct{})
	go func() {
		logEntriesIn, _ := eqfiles.ListenForLogs()
		go func() {
			for {
				logEntries := <-logEntriesIn
				for _, logEntry := range logEntries {
					c.parseForBid(ctx.EventEnv(), logEntry)
				}
			}
		}()
		bufferedLogs, err := rpcidl.FetchBufferedMessages(ipcrenderer.Client)
		if err != nil {
			console.Log("Failed to retreive buffered logs ", err)
		}
		for _, logEntry := range bufferedLogs {
			c.parseForBid(ctx.EventEnv(), logEntry)
		}
	}()
}

func extractNumbers(text string) []int {
	var result []int
	var buffer *strings.Builder
	for _, c := range text {
		if unicode.IsDigit(c) {
			if buffer == nil {
				buffer = new(strings.Builder)
			}
			buffer.WriteRune(c)
		} else {
			if buffer != nil {
				ival, _ := strconv.Atoi(buffer.String())
				result = append(result, ival)
				buffer = nil
			}
		}
	}
	if buffer != nil {
		ival, _ := strconv.Atoi(buffer.String())
		result = append(result, ival)
	}
	return result
}

func (c *Root) randomTiebreaker() int {
	for {
		rv := rand.Int() % len(commonWords)
		if _, ok := c.RandBag[rv]; !ok {
			c.RandBag[rv] = struct{}{}
			return rv
		}
	}
}

var tellRE = regexp.MustCompile("^([A-Za-z]+) (?:tells|told) you, '(.*)'$")

func (c *Root) parseForBid(env vugu.EventEnv, entry *eqfiles.LogEntry) {
	tellMatch := tellRE.FindStringSubmatch(entry.Message)
	if tellMatch == nil {
		return
	}
	sender := tellMatch[1]
	message := tellMatch[2]
	ivals := extractNumbers(message)
	updateOccurred := false

	for _, entry := range c.ActiveBids {
		if entry.Name == sender {
			if len(ivals) == 1 {
				entry.Value = ivals[0]
			}
			entry.Texts = append(entry.Texts, message)
			updateOccurred = true
			break
		}
	}
	if !updateOccurred {
		bid := &Bid{
			Name:       sender,
			Tiebreaker: c.randomTiebreaker(),
			Texts:      []string{message},
		}
		if len(ivals) == 1 {
			bid.Value = ivals[0]
		}
		c.ActiveBids = append(c.ActiveBids, bid)
	}

	env.Lock()
	sort.Sort(byValueDesc(c.ActiveBids))
	env.UnlockRender()
}

type byValueDesc []*Bid

func (b byValueDesc) Len() int { return len(b) }
func (b byValueDesc) Less(i, j int) bool {
	if b[i].Value == b[j].Value {
		return b[i].Tiebreaker < b[j].Tiebreaker
	} else {
		return b[i].Value > b[j].Value
	}
}
func (b byValueDesc) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
